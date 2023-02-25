> WIP

# 🔥 Forge

A Go library/tool for building fullstack webapps.

## Usage

Forge can be used either directly as a tool or as a library.

### As a tool

1. Download & unpack a Go release from [releases](https://github.com/spy16/forge/releases) section.
2. Tune the `forge.yml` file as per your needs.
3. Run `./forge serve -c forge.yml`

> If you have frontend build, use `--static=./ui` to serve as static files.

### As a library [Goal]

To use the pre-built CLI with custom initialisation of app:

```golang
package main

func main() {
	cli := forge.CLI("myapp",
		forge.WithFirebase(),
		forge.WithStatic(http.Dir("./foo")),
		forge.WithPostHook(func(app core.App, ge *gin.Engine) error {
			ge.GET("/api/myendpoint", app.Authenticate(), func(ctx *gin.Context) {
				// do some stuff
			})

			return nil
		}),
	)
	_ = cli.Execute()
}

```

To use Forge (or "forge an app") from scratch, use the `forge.Forge()` function.

```golang
package main

func main() {
	rawMaterials := []forge.Option{
		forge.WithConfLoader(myOwn),
		forge.WithPGBase(),
		forge.WithStatic(http.Dir("./foo")),
		forge.WithPostHook(func (app core.App, ge *gin.Engine) error {
            ge.GET("/api/myendpoint", app.Authenticate(), func(ctx *gin.Context) {
                // do some stuff
            })
			
			return nil
	    }),
	}

	ge, err := forge.Forge("myapp", rawMaterials...)
	if err != nil {
		panic(err)
	}
	ge.Run()
}
```
