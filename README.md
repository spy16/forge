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

### As a library

To use the pre-built CLI with custom initialisation of app:

```golang
package foo

func main() {
	postHook := func(app core.App, loader core.ConfLoader) error {
		router := app.Router()
		router.Get("/my-api", myHandler)
		return nil
	}

	cmd := forge.CLI("myapp",forge.WithPostHook(postHook))
	_ = cmd.Execute()
}

```

To use Forge (or "forge an app") from scratch, use the `forge.Forge()` function.

```golang
package foo

func main() {
	rawMaterials := []forge.Option{
		forge.WithAuth(customAuthModule),
		forge.WithConfLoader(customConfigLoader),
		forge.WithPostHook(func(app core.App, loader core.ConfLoader) error {
			router := app.Router()
			router.Get("/my-own-api", myHandler)
			return nil
		}),
		// ... more custom things if you want
	}

	app, _ := forge.Forge("myapp", rawMaterials...)
	_ = httpx.Serve(ctx, ":8080", router, 5*time.Second)
}
```