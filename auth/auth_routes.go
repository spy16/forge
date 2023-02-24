package auth

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/httpx"
	"github.com/spy16/forge/pkg/strutils"
)

const (
	redirectToParam = "redirect_to"

	contentTypeForm = "application/x-www-form-urlencoded"
)

// Routes installs auth module routes onto the given router.
func (auth *Auth) Routes(r chi.Router) {
	r.Post("/register", auth.handleRegister)
	r.Post("/login", auth.handleLogin)
	r.Get("/logout", auth.handleLogout)

	r.Get("/oauth2", auth.handleOAuth2Redirect)
	r.Get("/oauth2/cb", auth.handleOAuth2Callback)

	r.Group(func(r chi.Router) {
		r.Use(auth.Authenticate())

		r.Get("/me", httpx.WrapErrH(auth.handleWhoAmI))
	})
}

func (auth *Auth) handleRegister(w http.ResponseWriter, r *http.Request) {
	doRegister := func() (*User, error) {
		var creds userCreds
		if err := creds.readFrom(r); err != nil {
			return nil, err
		} else if !strutils.IsValidEmail(creds.Email) {
			return nil, errors.MissingAuth.Hintf("invalid email")
		}

		pwdHash, err := HashPassword(creds.Password)
		if err != nil {
			return nil, err
		}

		u := NewUser(creds.Kind, creds.Username, creds.Email)
		u.PwdHash = &pwdHash

		registeredU, err := auth.RegisterUser(r.Context(), u, nil)
		if err != nil {
			if errors.OneOf(err, []error{errors.Conflict, errors.InvalidInput}) {
				return nil, err
			}
			return nil, errors.InternalIssue.CausedBy(err)
		}
		return registeredU, nil
	}

	u, err := doRegister()
	if err != nil {
		writeErr(w, r, auth.cfg.RegisterPageRoute, err)
	} else {
		writeSuccess(w, r, auth.cfg.RegisterPageRoute, http.StatusCreated, u.Clone(true))
	}
}

func (auth *Auth) handleOAuth2Redirect(w http.ResponseWriter, r *http.Request) {
	prepareRedirection := func() (string, *oauth2FlowState, error) {
		q := r.URL.Query()
		userKind := q.Get("kind")
		providerID := q.Get("p")
		if userKind == "" {
			userKind = defaultUserKind
		}
		if !strutils.OneOf(userKind, auth.cfg.EnabledKinds) {
			return "", nil, errors.InvalidInput.Coded("invalid_kind").Hintf("user kind '%s' is not valid", userKind)
		}

		p, err := goth.GetProvider(providerID)
		if err != nil {
			return "", nil, errors.InvalidInput.Coded("invalid_provider").CausedBy(err)
		}
		state := strutils.RandStr(10)

		sess, err := p.BeginAuth(state)
		if err != nil {
			return "", nil, errors.InternalIssue.CausedBy(err)
		}

		authURL, err := sess.GetAuthURL()
		if err != nil {
			return "", nil, errors.InternalIssue.CausedBy(err)
		}

		return authURL, &oauth2FlowState{
			UserKind:   userKind,
			Provider:   p.Name(),
			Session:    sess.Marshal(),
			RedirectTo: r.FormValue(redirectToParam),
		}, nil
	}

	authURL, state, err := prepareRedirection()
	if err != nil {
		writeErr(w, r, auth.cfg.LoginPageRoute, err)
		return
	}

	setOAuthFlowState(w, state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (auth *Auth) handleOAuth2Callback(w http.ResponseWriter, r *http.Request) {
	var errInvalidCB = errors.InvalidInput.Coded("invalid_callback")

	processCallback := func() (*User, error) {
		flowState := popOAuthState(w, r)
		if flowState == nil {
			return nil, errInvalidCB.Hintf("oauth2 flow state is nil")
		}

		p, err := goth.GetProvider(flowState.Provider)
		if err != nil {
			return nil, errInvalidCB.CausedBy(err)
		}

		sess, err := p.UnmarshalSession(flowState.Session)
		if err != nil {
			return nil, errInvalidCB.CausedBy(err)
		}

		q := r.URL.Query()
		if !checkCallbackState(sess, q.Get("state")) {
			return nil, errInvalidCB.Hintf("state value mismatch")
		}

		if _, err := sess.Authorize(p, q); err != nil {
			return nil, errors.InternalIssue.CausedBy(err)
		}

		gothUser, err := p.FetchUser(sess)
		if err != nil {
			return nil, errors.InternalIssue.CausedBy(err)
		}

		loginKeyID := NewAuthKey(gothUser.Provider, gothUser.UserID)
		exU, err := auth.GetUser(r.Context(), loginKeyID)
		if err != nil && !errors.Is(err, errors.NotFound) {
			return nil, errors.InternalIssue.CausedBy(err)
		}

		if exU == nil {
			// new user registration
			newU := NewUser(flowState.UserKind, "", gothUser.Email)
			newU.Data = userDataFromGothUser(gothUser)

			loginKey := Key{
				Key: loginKeyID,
				Attribs: map[string]any{
					"user_id":       gothUser.UserID,
					"expires_at":    gothUser.ExpiresAt.Unix(),
					"access_token":  gothUser.AccessToken,
					"refresh_token": gothUser.RefreshToken,
					"raw_data":      gothUser.RawData,
				},
			}

			exU, err = auth.RegisterUser(r.Context(), newU, []Key{loginKey})
			if err != nil {
				if !errors.OneOf(err, []error{errors.Conflict}) {
					err = errors.InternalIssue.CausedBy(err)
				}
				return nil, err
			}
		} else {
			// TODO: update existing user
		}
		return exU, nil
	}

	u, err := processCallback()
	if err != nil {
		writeErr(w, r, auth.cfg.LoginPageRoute, err)
		return
	}
	auth.finishLogin(w, r, *u)
}

func (auth *Auth) handleVerify(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	id := q.Get("id")
	verifyToken := q.Get("token")

	u, err := auth.VerifyUser(r.Context(), id, verifyToken)
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			err = errors.MissingAuth
		}
		writeErr(w, r, auth.cfg.LoginPageRoute, err)
		return
	}

	auth.finishLogin(w, r, *u)
}

func (auth *Auth) handleLogin(w http.ResponseWriter, r *http.Request) {
	doLogin := func() (*User, error) {
		var creds userCreds
		if err := creds.readFrom(r); err != nil {
			return nil, err
		}

		keyKind := KeyKindUsername
		keyValue := creds.Username
		if creds.Email != "" {
			keyKind = KeyKindEmail
			keyValue = creds.Email
		}

		u, err := auth.GetUser(r.Context(), NewAuthKey(keyKind, keyValue))
		if err != nil {
			if errors.Is(err, errors.NotFound) {
				err = errors.MissingAuth.Hintf("user not found")
			}
			return nil, err
		} else if !u.CheckPassword(creds.Password) {
			return nil, errors.MissingAuth.Hintf("password mismatch")
		} else if creds.Email != "" && creds.Email != u.Email {
			return nil, errors.MissingAuth.Hintf("email mismatch")
		} else if u.Kind != creds.Kind {
			return nil, errors.MissingAuth.Hintf("user kind mismatch")
		}

		return u, nil
	}

	u, err := doLogin()
	if err != nil {
		writeErr(w, r, auth.cfg.LoginPageRoute, err)
		return
	}

	auth.finishLogin(w, r, *u)
}

func (auth *Auth) handleWhoAmI(w http.ResponseWriter, r *http.Request) error {
	session := CurSession(r.Context())
	if session == nil {
		return errors.MissingAuth
	}

	u, err := auth.GetUser(r.Context(), NewAuthKey(KeyKindID, session.UserID))
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			return errors.MissingAuth
		}
		return errors.InternalIssue.CausedBy(err)
	}

	httpx.WriteJSON(w, r, http.StatusOK, u.Clone(true))
	return nil
}

func (auth *Auth) handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     auth.cfg.SessionCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (auth *Auth) finishLogin(w http.ResponseWriter, r *http.Request, user User) {
	session, err := auth.CreateSession(r.Context(), user)
	if err != nil {
		writeErr(w, r, auth.cfg.LoginPageRoute, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     auth.cfg.SessionCookie,
		Value:    session.Token,
		Path:     "/",
		Expires:  session.ExpiresAt,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeSuccess(w, r, auth.cfg.LoginPageRoute, http.StatusOK, map[string]any{
		"user":   user.Clone(true),
		"token":  session.Token,
		"expiry": session.ExpiresAt,
	})
}

func userDataFromGothUser(gu goth.User) UserData {
	// TODO: extract more from raw data?
	return map[string]any{
		"name":      gu.Name,
		"picture":   gu.AvatarURL,
		"location":  gu.Location,
		"nick_name": gu.NickName,
	}
}

func writeErr(w http.ResponseWriter, r *http.Request, redirectTo string, err error) {
	isFormSubmit := strings.Contains(r.Header.Get("Content-Type"), contentTypeForm)
	if isFormSubmit {
		u, parseErr := url.Parse(redirectTo)
		if parseErr == nil {
			q := u.Query()
			q.Set("err_code", errors.E(err).Code)
			u.RawQuery = q.Encode()
			u.JoinPath()
			http.Redirect(w, r, u.String(), http.StatusSeeOther)
			return
		}
	}
	httpx.WriteErr(w, r, err)
}

func writeSuccess(w http.ResponseWriter, r *http.Request, redirectTo string, status int, v any) {
	isFormSubmit := strings.Contains(r.Header.Get("Content-Type"), contentTypeForm)
	if isFormSubmit {
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
	} else {
		httpx.WriteJSON(w, r, status, v)
	}
}
