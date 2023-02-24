package router

import (
	"net/http"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/pkg/httpx"
)

func handleWhoAmI(app core.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rc := core.FromCtx(r.Context())
		if !rc.Authenticated() {
			httpx.WriteErr(w, r, errors.MissingAuth)
			return
		}

		u, err := app.Substrate().User(r.Context(), rc.Session.UserID)
		if err != nil {
			if errors.Is(err, errors.NotFound) {
				err = errors.MissingAuth
			} else {
				err = errors.InternalIssue.CausedBy(err)
			}
			httpx.WriteErr(w, r, err)
			return
		}

		httpx.WriteJSON(w, r, http.StatusOK, u.Clone(true))
	}
}
