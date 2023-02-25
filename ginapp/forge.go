package ginapp

import (
	"context"
	"time"

	"github.com/kisielk/gotool"

	"github.com/spy16/forge/core"
)

type App struct {
	auth  Auth
	oauth OAuth2Method
	creds CredsMethod
}

type Auth interface {
	Authenticate(ctx gotool.Context, token string) (*Session, error)
}

type OAuth2Method interface {
	AuthURL(ctx context.Context, provider string) (*FlowState, error)
	CodeLogin(ctx context.Context, code, state string, fs FlowState) (*Session, error)
}

type CredsMethod interface {
	Register(ctx context.Context, creds UserCreds) (*core.User, error)
	PwdLogin(ctx context.Context, creds UserCreds) (*Session, error)
}

type Session struct {
	User   core.User
	Token  string
	Expiry time.Time
}

type UserCreds struct {
	Email    string `json:"email" form:"email"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type FlowState struct {
	State       string `json:"state"`
	Provider    string `json:"provider"`
	RedirectURL string `json:"redirect_url"`
}
