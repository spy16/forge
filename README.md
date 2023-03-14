> WIP

# 🔥 Forge

A Go library/tool for building backend for fullstack apps.

## Usage

Forge can be used either directly as a tool or as a library.

### As a tool

1. Download & unpack a Go release from [releases](https://github.com/spy16/forge/releases) section.
2. Tune the `forge.yml` file as per your needs.
3. Run `./forge serve -c forge.yml`

> If you have frontend build, use `--static=./ui` to serve as static files.

### As a library [Goal]

```golang
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spy16/forge"
	"github.com/spy16/forge/builtins/firebase"
	"github.com/spy16/forge/core"
)

func main() {
	cli := forge.CLI("myapp",
		forge.WithAuth(&firebase.Auth{
			ProjectID: "foo",
		}),
		forge.WithPostHook(func(app core.App, conf core.ConfigLoader) error {
			r := app.Chi()
			r.Use(app.Authenticate())
			r.Get("/api/my-endpoint", func(w http.ResponseWriter, r *http.Request) {
				// Only accessible with firebase auth token
            })
			return nil
		}),
	)
	_ = cli.Execute()
}
```
