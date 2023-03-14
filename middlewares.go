package forge

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/log"
	"github.com/spy16/forge/core/servio"
)

func extractReqCtx() core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rc := core.ReqCtx{
				Path:       r.URL.Path,
				Route:      r.URL.Path, // TODO: set route pattern here.
				Method:     r.Method,
				Session:    nil,
				RequestID:  middleware.GetReqID(r.Context()),
				RemoteAddr: r.RemoteAddr,
			}

			ctx := core.NewCtx(r.Context(), rc)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func requestLogger() core.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now()

			rwc := &servio.ResponseWriterCapture{ResponseWriter: w}
			next.ServeHTTP(rwc, r)

			status := rwc.Status
			fields := core.M{
				"status":  status,
				"latency": time.Since(t),
			}

			if status >= 500 {
				log.Error(r.Context(), "request finished with 5xx", nil, fields)
			} else if status >= 400 {
				log.Warn(r.Context(), "request finished with 4xx", fields)
			} else {
				log.Info(r.Context(), "request finished", fields)
			}
		})
	}
}
