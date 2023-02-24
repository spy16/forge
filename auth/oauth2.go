package auth

import (
	"net/http"
	"net/url"
	"time"

	"github.com/markbates/goth"

	"github.com/spy16/forge/pkg/httpx"
)

const (
	oauthFlowCookie = "_oauth_state"
	oauthCookieTTL  = 10 * time.Minute
)

type oauth2FlowState struct {
	UserKind   string `json:"user_kind"`
	Provider   string `json:"provider"`
	Session    string `json:"goth_session"`
	RedirectTo string `json:"redirect_to"`
}

func setOAuthFlowState(w http.ResponseWriter, state *oauth2FlowState) {
	now := time.Now()
	c := &http.Cookie{
		Name:     oauthFlowCookie,
		Value:    "",
		Path:     "/",
		Expires:  now,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	if state != nil {
		value, err := httpx.MarshalCookie(state)
		if err != nil {
			panic(err)
		}
		c.Value = value
		c.Expires = now.Add(oauthCookieTTL)
	}

	http.SetCookie(w, c)
}

func popOAuthState(w http.ResponseWriter, r *http.Request) *oauth2FlowState {
	defer setOAuthFlowState(w, nil)

	var st oauth2FlowState
	if !httpx.UnmarshalCookie(r, oauthFlowCookie, &st) {
		return nil
	}
	return &st
}

func checkCallbackState(sess goth.Session, actualState string) bool {
	authURL, err := sess.GetAuthURL()
	if err != nil {
		return false
	}

	u, err := url.Parse(authURL)
	if err != nil {
		return false
	}

	return u.Query().Get("state") == actualState
}
