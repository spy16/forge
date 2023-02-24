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

type wrappedRW struct {
	http.ResponseWriter
	status int
}

func (w *wrappedRW) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func reqLog() Middleware {
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

			wr := &wrappedRW{
				status:         http.StatusOK,
				ResponseWriter: w,
			}
			next.ServeHTTP(wr, r.WithContext(ctx))

			log.Info(ctx, "request finished", core.M{
				"status":  wr.status,
				"latency": time.Since(t),
			})
		})
	}
}

func extractReqCtx(auth core.Auth, cookieName string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			rc := core.ReqCtx{
				Path:       r.URL.Path,
				Method:     r.Method,
				Session:    nil,
				RequestID:  middleware.GetReqID(ctx),
				RemoteAddr: r.RemoteAddr,
			}

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
				rc.Session = sess
			}

			next.ServeHTTP(w, r.WithContext(core.NewCtx(ctx, rc)))
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
