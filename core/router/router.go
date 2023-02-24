package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/spy16/forge/core"
)

const (
	bearerPrefix = "Bearer "
	headerAuthz  = "Authorization"

	defAuthCookie  = "_forge_auth"
	defRoutePrefix = "/forge"
)

// New returns a new API router instance for the given app.
func New(app core.App, cnfL core.ConfLoader) (chi.Router, error) {
	prefix := cnfL.String("forge.router.prefix", defRoutePrefix)
	cookieName := cnfL.String("forge.auth.cookie_name", defAuthCookie)

	router := chi.NewRouter()
	router.Use(
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
		extractReqCtx(),
		extractLogCtx(),
		authenticate(app.Auth(), cookieName),
	)

	router.Route(prefix, func(r chi.Router) {
		r.Get("/ping", handlePing())
		r.Get("/me", handleWhoAmI(app))
	})

	return router, nil
}

func handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}
}
