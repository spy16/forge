package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/log"
	"github.com/spy16/forge/pkg/httpx"
)

type Middleware func(http.Handler) http.Handler

func extractLogCtx() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()
			rc := core.FromCtx(r.Context())
			ctx := log.Ctx(r.Context(), core.M{
				"path":        r.URL.Path,
				"method":      r.Method,
				"req_id":      rc.RequestID,
				"authn":       rc.Authenticated(),
				"remote_addr": rc.RemoteAddr,
			})

			wr := &httpx.WrappedResponseWriter{
				Status:         http.StatusOK,
				ResponseWriter: w,
			}
			next.ServeHTTP(wr, r.WithContext(ctx))

			if wr.Status >= 500 {
				log.Error(ctx, "request finished with 5xx", wr.Error, core.M{
					"status":  wr.Status,
					"latency": time.Since(t),
				})
			} else if wr.Status >= 400 {
				log.Warn(ctx, "request finished with 4xx", core.M{
					"error":   wr.Error,
					"status":  wr.Status,
					"latency": time.Since(t),
				})
			} else {
				log.Info(ctx, "request finished", core.M{
					"status":  wr.Status,
					"latency": time.Since(t),
				})
			}
		})
	}
}

func authenticate(auth core.Auth, cookieName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// extract token and restore session if any.
			if token := extractToken(r, cookieName); token != "" {
				sess, err := auth.RestoreSession(r.Context(), token)
				if err != nil {
					if !errors.Is(err, errors.MissingAuth) {
						err = errors.InternalIssue.CausedBy(err)
					}
					httpx.WriteErr(w, r, err)
					return
				}

				ctx := r.Context()
				rc := core.FromCtx(ctx)
				rc.Session = sess
				ctx = core.NewCtx(ctx, rc)

				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func extractReqCtx() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = core.NewCtx(r.Context(), core.ReqCtx{
				Path:       r.URL.Path,
				Method:     r.Method,
				Session:    nil,
				RequestID:  middleware.GetReqID(ctx),
				RemoteAddr: r.RemoteAddr,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request, cookieName string) string {
	var token string
	if authH := r.Header.Get(headerAuthz); strings.HasPrefix(authH, bearerPrefix) {
		token = strings.TrimPrefix(authH, bearerPrefix)
	} else {
		c, err := r.Cookie(cookieName)
		if err != nil || c == nil {
			return ""
		}
		token = c.Value
	}
	return strings.TrimSpace(token)
}
