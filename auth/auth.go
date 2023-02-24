package auth

import (
	"context"
	_ "embed"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/spy16/forge/core/errors"
)

const defaultSessionCookie = "_forge_auth"

//go:embed schema.sql
var schema string

// Init initialises auth module and returns.
func Init(conn *pgx.Conn, baseURL string, cfg Config) (*Auth, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, errors.InvalidInput.Hintf("invalid baseURL").CausedBy(err)
	}

	if err := cfg.sanitise(u); err != nil {
		return nil, err
	}

	cbURL := u.JoinPath("/oauth2/cb").String()
	goth.UseProviders(
		google.New(cfg.Google.ClientID, cfg.Google.ClientSecret, cbURL, cfg.Google.Scopes...),
		github.New(cfg.Github.ClientID, cfg.Github.ClientSecret, cbURL, cfg.Github.Scopes...),
	)

	if _, err := conn.Exec(context.Background(), schema); err != nil {
		return nil, err
	}

	au := &Auth{
		cfg:  cfg,
		conn: conn,
		providers: map[string]goth.Provider{
			"google": google.New(cfg.Google.ClientID, cfg.Google.ClientSecret, cbURL, cfg.Google.Scopes...),
			"github": github.New(cfg.Github.ClientID, cfg.Github.ClientSecret, cbURL, cfg.Github.Scopes...),
		},
	}

	return au, nil
}

// Auth represents the auth module and implements user management and
// authentication facilities.
type Auth struct {
	cfg       Config
	conn      *pgx.Conn
	providers map[string]goth.Provider
}

type Config struct {
	SessionTTL    time.Duration `mapstructure:"session_ttl"`
	SessionCookie string        `mapstructure:"session_cookie"`
	SigningSecret string        `mapstructure:"signing_secret"`
	EnabledKinds  []string      `mapstructure:"enabled_kinds"`

	LoginPageRoute    string `mapstructure:"login_page_route"`
	RegisterPageRoute string `mapstructure:"register_page_route"`

	Google OAuthConf `mapstructure:"google"`
	Github OAuthConf `mapstructure:"github"`
}

type OAuthConf struct {
	Scopes       []string `mapstructure:"scopes"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
}

func (cfg *Config) sanitise(u *url.URL) error {
	if cfg.RegisterPageRoute != "" {
		cfg.RegisterPageRoute = u.JoinPath(cfg.RegisterPageRoute).String()
	}

	if cfg.LoginPageRoute != "" {
		cfg.LoginPageRoute = u.JoinPath(cfg.LoginPageRoute).String()
	}

	if cfg.SessionTTL <= 0 {
		cfg.SessionTTL = 12 * time.Hour
	}

	if cfg.SessionCookie == "" {
		cfg.SessionCookie = defaultSessionCookie
	}

	if cfg.SigningSecret == "" {
		return errors.InvalidInput.Hintf("signing_secret is required")
	}

	if len(cfg.EnabledKinds) == 0 {
		cfg.EnabledKinds = []string{defaultUserKind}
	}

	return nil
}
