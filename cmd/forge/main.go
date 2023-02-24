package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spy16/forge"
)

var (
	Commit    = "n/a"
	Version   = "n/a"
	BuildTime = "n/a"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cmd := forge.CLI("forge")
	cmd.Short = "🔥 A tiny Go platform for building SaaS applications."
	cmd.Version = fmt.Sprintf("%s\ncommit: %s\nbuilt_at: %s", Version, Commit, BuildTime)
	cmd.SetContext(ctx)
	_ = cmd.Execute()
}
