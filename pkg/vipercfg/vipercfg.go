package vipercfg

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Init initialises the config loader with given options and returns.
func Init(opts ...Option) (*Loader, error) {
	l := &Loader{viper: viper.New()}
	for _, opt := range withDefault(opts) {
		opt(l)
	}
	l.viper.SetConfigName(l.confName)

	// for transforming app.host to app_host
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	if l.envPrefix != "" {
		l.viper.SetEnvPrefix(l.envPrefix)
	}

	if l.confFile != "" {
		l.viper.SetConfigFile(l.confFile)
		if err := l.viper.ReadInConfig(); err != nil {
			return nil, err
		}
	} else {
		for _, dir := range l.confDirs {
			l.viper.AddConfigPath(dir)
		}
		_ = l.viper.ReadInConfig()
	}

	l.viper.AutomaticEnv()
	return l, nil
}

// Loader implements config loader facilities using Viper.
type Loader struct {
	viper     *viper.Viper
	confFile  string
	confDirs  []string
	confName  string
	envPrefix string
}

func (l *Loader) Viper() *viper.Viper { return l.viper }

// Int returns the int value set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) Int(key string, defaultValue int) int {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetInt(key)
}

// Bool returns the boolean value set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) Bool(key string, defaultValue bool) bool {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetBool(key)
}

// String returns the string value set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) String(key string, defaultValue string) string {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetString(key)
}

// Strings returns the list of string values set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) Strings(key string, defaultValue []string) []string {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetStringSlice(key)
}

// Float64 returns the float64 value set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) Float64(key string, defaultValue float64) float64 {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetFloat64(key)
}

// Duration returns the duration value set for the given key.
// Returns defaultValue if keys is not explicitly set.
func (l *Loader) Duration(key string, defaultValue time.Duration) time.Duration {
	if !l.viper.IsSet(key) {
		return defaultValue
	}
	return l.viper.GetDuration(key)
}

func getExecPath() string {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		baseDir, _ := os.Getwd()
		return baseDir
	} else {
		return filepath.Dir(os.Args[0])
	}
}
