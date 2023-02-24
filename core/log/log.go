package log

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/spy16/forge/core"
	"github.com/spy16/forge/core/errors"
)

var lg = logrus.New()

// Setup configures the global logger with given log-level
// and output format. JSON and text formats are supported.
func Setup(level, format string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.WarnLevel
	}

	lg.SetLevel(lvl)
	lg.SetFormatter(&logrus.TextFormatter{})
	if format == "json" {
		lg.SetFormatter(&logrus.JSONFormatter{})
	}
}

func Debug(ctx context.Context, msg string, fields ...core.M) {
	doLog(ctx, logrus.DebugLevel, msg, nil, fields)
}

func Info(ctx context.Context, msg string, fields ...core.M) {
	doLog(ctx, logrus.InfoLevel, msg, nil, fields)
}

func Warn(ctx context.Context, msg string, fields ...core.M) {
	doLog(ctx, logrus.WarnLevel, msg, nil, fields)
}

func Error(ctx context.Context, msg string, err error, fields ...core.M) {
	doLog(ctx, logrus.ErrorLevel, msg, err, fields)
}

func Fatal(ctx context.Context, msg string, err error, fields ...core.M) {
	doLog(ctx, logrus.FatalLevel, msg, err, fields)
	os.Exit(1)
}

func doLog(ctx context.Context, level logrus.Level, msg string, err error, args []core.M) {
	if ctx == nil {
		ctx = context.Background()
	}

	m := mergeFields(fromCtx(ctx), args)
	if err != nil {
		e := errors.E(err)
		m["error"] = err.Error()
		if len(e.Attribs) > 0 {
			m["error.attribs"] = e.Attribs
		}
	}
	lg.WithFields(m).Logln(level, msg)
}

func mergeFields(base core.M, fields []core.M) map[string]any {
	res := map[string]any{}
	for k, v := range base {
		res[k] = v
	}
	for _, fieldMap := range fields {
		for k, v := range fieldMap {
			res[k] = v
		}
	}
	return res
}
