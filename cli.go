package forge

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/log"
	"github.com/spy16/forge/core/servio"
	"github.com/spy16/forge/core/vipercfg"
)

// CLI returns a new Cobra CLI that can be used directly.
func CLI(name string, forgeOpts ...Option) *cobra.Command {
	if !namePattern.MatchString(name) {
		panic(errInvalidName)
	}

	cli := &cobra.Command{
		Use:   fmt.Sprintf("%s <command> [flags] [args]", name),
		Short: fmt.Sprintf("%s: a forge application", name),
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	flags := cli.PersistentFlags()
	flags.StringP("config", "c", "", "Override config file path")
	flags.StringP("log-level", "L", "info", "Min log level to start from")
	flags.String("log-format", "text", "Log output format (json or text)")

	cli.AddCommand(
		cmdServe(name, forgeOpts),
		cmdConfigs(name),
	)
	return cli
}

func cmdServe(name string, forgeOpts []Option) *cobra.Command {
	var httpAddr, staticDir string
	var graceT time.Duration
	var debug bool
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start HTTP server",
		Run: func(cmd *cobra.Command, args []string) {
			cl := makeConfLoader(name, cmd)
			forgeOpts = append(forgeOpts,
				WithConfLoader(cl),
				WithPreHook(func(app PreContext) error {
					// TODO: initialise modules using config only.
					return nil
				}),
			)

			app, err := Forge(name, forgeOpts...)
			if err != nil {
				log.Fatal(cmd.Context(), "failed to forge app", err)
			}

			if staticDir != "" {
				app.Mount("/", http.FileServer(http.Dir(staticDir)))
			}

			log.Info(cmd.Context(), "starting http server", core.M{"http_addr": httpAddr})
			if err := servio.Serve(cmd.Context(), httpAddr, app); err != nil {
				log.Fatal(cmd.Context(), "server exited with error", err)
			}
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&httpAddr, "http", ":8080", "HTTP server address")
	flags.StringVar(&staticDir, "static", "", "If set, serves all files in the dir as-is")
	flags.DurationVarP(&graceT, "grace", "G", 5*time.Second, "Grace period for shutdown")
	flags.BoolVar(&debug, "debug", false, "Running in debug mode")

	return cmd
}

func cmdConfigs(name string) *cobra.Command {
	return &cobra.Command{
		Use: "configs",
		Run: func(cmd *cobra.Command, args []string) {
			cnfL := makeConfLoader(name, cmd)
			v := cnfL.Viper()

			var buf bytes.Buffer
			_ = yaml.NewEncoder(&buf).Encode(v.AllSettings())

			fmt.Printf("# file: %s\n\n%s", v.ConfigFileUsed(), buf.String())
		},
	}
}

func makeConfLoader(name string, cmd *cobra.Command) *vipercfg.Loader {
	cl, err := vipercfg.Init(
		vipercfg.WithName(name),
		vipercfg.WithCobra(cmd, "config"),
	)
	if err != nil {
		log.Fatal(cmd.Context(), "failed to load configs", err)
	}

	log.Setup(
		cl.String("log_level", "info"),
		cl.String("log_format", "text"),
	)
	return cl
}
