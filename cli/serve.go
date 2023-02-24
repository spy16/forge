package cli

import (
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/spy16/forge"
	"github.com/spy16/forge/core"
	"github.com/spy16/forge/pkg/httpx"
	"github.com/spy16/forge/pkg/log"
)

func cmdServe(name string, hooks []Hook) *cobra.Command {
	var httpAddr, staticDir string
	var graceT time.Duration
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			cl := makeConfLoader(name, cmd)

			app, err := forge.Forge(cmd.Context(), name, cl)
			if err != nil {
				log.Fatal(cmd.Context(), "failed to forge app", err)
			}

			for _, hook := range hooks {
				if err := hook(app); err != nil {
					log.Fatal(cmd.Context(), "failed to forge app", err)
				}
			}

			router := app.Router()
			if staticDir != "" {
				router.Mount("/", http.FileServer(http.Dir(staticDir)))
			}

			log.Info(cmd.Context(), "starting http server", core.M{"http_addr": httpAddr})
			if err := httpx.Serve(cmd.Context(), httpAddr, router, graceT); err != nil {
				log.Fatal(cmd.Context(), "server exited with error", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&httpAddr, "http", ":8080", "HTTP server address")
	flags.StringVar(&staticDir, "static", "", "If set, serves all files in the dir as-is")
	flags.DurationVarP(&graceT, "grace", "G", 5*time.Second, "Grace period for shutdown")
	return cmd
}
