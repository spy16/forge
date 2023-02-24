package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/pkg/log"
	"github.com/spy16/forge/pkg/vipercfg"
)

type Hook func(core.App) error

// New returns a new Cobra CLI that can be used directly.
func New(name string, hooks ...Hook) *cobra.Command {
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
		cmdServe(name, hooks),
	)
	return cli
}

func makeConfLoader(name string, cmd *cobra.Command) *vipercfg.Loader {
	cl, err := vipercfg.Init(
		vipercfg.WithName(name),
		vipercfg.WithCobra(cmd, "config"),
	)
	if err != nil {
		log.Fatal(cmd.Context(), "failed to load configs", err)
	}
	return cl
}
