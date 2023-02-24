package vipercfg

import (
	"strings"

	"github.com/spf13/cobra"
)

// Option values can be provided to Load() for customisation.
type Option func(l *Loader)

// WithEnvPrefix sets a prefix to be used for all config variables
// from environment.
func WithEnvPrefix(prefix string) Option {
	return func(l *Loader) {
		l.envPrefix = strings.TrimSpace(prefix)
	}
}

// WithName sets the configuration file title to be used for automatic
// discovery of config files.
func WithName(name string) Option {
	return func(l *Loader) {
		l.confName = strings.TrimSpace(name)
	}
}

// WithCobra enables reading config file overrides from a flag. When the
// flag is specified, it acts as the only file-based source of configs.
// If not specified, config files are auto-discovered.
func WithCobra(cmd *cobra.Command, flagName string) Option {
	return func(l *Loader) {
		cfgFile, _ := cmd.Flags().GetString(flagName)
		if cfgFile != "" {
			l.confFile = cfgFile
		}
	}
}

// WithFile sets a config file to use explicitly. When the filePath is
// not empty, it acts the only file-based source of configs. If empty,
// config files are auto-discovered.
func WithFile(filePath string) Option {
	return func(l *Loader) {
		l.confFile = filePath
	}
}

// WithPaths overrides the default directories that are searched for
// config files.
func WithPaths(paths ...string) Option {
	return func(l *Loader) {
		l.confDirs = paths
	}
}

func withDefault(opts []Option) []Option {
	return append([]Option{
		WithName("config"),
		WithPaths("./", getExecPath()),
	}, opts...)
}
