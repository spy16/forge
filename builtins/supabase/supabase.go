package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
)

// Auth implements auth module for supabase-based user management.
type Auth struct {
	APIKey    string `json:"api_key"`
	ProjectID string `json:"project_id"`
}

func (sb *Auth) Authenticate(ctx context.Context, token string) (*core.Session, error) {
	userURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/user", sb.ProjectID)

	req, err := http.NewRequest(http.MethodGet, userURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"APIKey":        {sb.APIKey},
		"Authorization": {"Bearer " + token},
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == 401 {
		return nil, errors.MissingAuth.Hintf("Auth returned 401")
	} else if resp.StatusCode != 200 {
		return nil, errors.InternalIssue.Hintf("Auth returned unexpected status: %s", resp.Status)
	}

	var userData supabaseUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return nil, errors.InternalIssue.Hintf("Auth returned invalid response").CausedBy(err)
	}

	return &core.Session{
		User: core.User{
			ID: userData.ID,
			Data: map[string]any{
				"name":    userData.UserMetadata.Name,
				"picture": userData.UserMetadata.AvatarURL,
			},
			Email:     userData.Email,
			CreatedAt: userData.CreatedAt,
			UpdatedAt: userData.UpdatedAt,
		},
		Token: token,
	}, nil
}

type supabaseUserResponse struct {
	ID               string    `json:"id"`
	Aud              string    `json:"aud"`
	Role             string    `json:"role"`
	Email            string    `json:"email"`
	EmailConfirmedAt time.Time `json:"email_confirmed_at"`
	Phone            string    `json:"phone"`
	ConfirmedAt      time.Time `json:"confirmed_at"`
	RecoverySentAt   time.Time `json:"recovery_sent_at"`
	LastSignInAt     time.Time `json:"last_sign_in_at"`
	AppMetadata      struct {
		Provider  string   `json:"provider"`
		Providers []string `json:"providers"`
	} `json:"app_metadata"`
	UserMetadata struct {
		AvatarURL     string `json:"avatar_url"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		FullName      string `json:"full_name"`
		Iss           string `json:"iss"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		ProviderID    string `json:"provider_id"`
		Sub           string `json:"sub"`
	} `json:"user_metadata"`
	Identities []struct {
		ID           string `json:"id"`
		UserID       string `json:"user_id"`
		IdentityData struct {
			AvatarURL     string `json:"avatar_url"`
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			FullName      string `json:"full_name"`
			Iss           string `json:"iss"`
			Name          string `json:"name"`
			Picture       string `json:"picture"`
			ProviderID    string `json:"provider_id"`
			Sub           string `json:"sub"`
		} `json:"identity_data"`
		Provider     string    `json:"provider"`
		LastSignInAt time.Time `json:"last_sign_in_at"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	} `json:"identities"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
