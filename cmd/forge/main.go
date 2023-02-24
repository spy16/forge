package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spy16/forge/cli"
)

var (
	Commit    = "n/a"
	Version   = "n/a"
	BuildTime = "n/a"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cmd := cli.New("forge")
	cmd.SetContext(ctx)
	cmd.Version = fmt.Sprintf("%s\ncommit: %s\nbuilt_at: %s", Version, Commit, BuildTime)
	_ = cmd.Execute()
}
