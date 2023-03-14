package forge

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/servio"
)

const defRoutePrefix = "/forge"

func (app *forgedApp) setupRoutes() error {
	app.chi.Route(defRoutePrefix, func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			servio.JSON(w, r, http.StatusNoContent, nil)
		})

		// authenticated routes.
		r.Group(func(r chi.Router) {
			r.Use(app.Authenticate())

			r.Get("/me", func(w http.ResponseWriter, r *http.Request) {
				rc := core.FromCtx(r.Context())
				servio.JSON(w, r, 200, rc.Session.User)
			})
		})
	})
	return nil
}
