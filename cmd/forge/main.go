package main

import (
	"fmt"

	"github.com/spy16/forge"
)

var (
	Commit    = "n/a"
	Version   = "n/a"
	BuildTime = "n/a"
)

func main() {
	cmd := forge.CLI("forge")
	cmd.Short = "🔥 A tiny Go platform for building SaaS applications."
	cmd.Version = fmt.Sprintf("%s\ncommit: %s\nbuilt_at: %s", Version, Commit, BuildTime)
	_ = cmd.Execute()
}
