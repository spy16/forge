package ginapp

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"

	"github.com/spy16/forge/core"
)

func New() (*gin.Engine, error) {
	ga := &ginApp{gin: gin.New()}
	if err := ga.setup(); err != nil {
		return nil, err
	}
	return ga.gin, nil
}

type ginApp struct {
	gin  *gin.Engine
	opts options
}

type options struct {
	authEnabled  bool
	registration bool
}

func (ga *ginApp) setup() error {
	ga.gin.Use(
		requestid.New(),
		extractReqCtx(),
		requestLogger(),
	)

	api := ga.gin.Group("/forge")

	api.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	authGrp := api.Group("/").Use(Authenticate())
	if ga.opts.authEnabled {
		api.POST("/logout", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, core.User{})
		})

		authGrp.GET("/me", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, core.User{})
		})
	}

	if ga.opts.registration {
		authGrp.POST("/register", func(ctx *gin.Context) {
			// TODO: username/email + password based registration
		})

		authGrp.POST("/login", func(ctx *gin.Context) {
			// TODO: username/email + password based registration
		})

		authGrp.GET("/oauth2", func(ctx *gin.Context) {
			// TODO: oauth2 redirection
		})
	}

	return nil
}

func extractToken(c *gin.Context) string {
	var token string
	const bearerPrefix = "Bearer "
	if authH := c.GetHeader("Authorization"); strings.HasPrefix(authH, bearerPrefix) {
		return strings.TrimPrefix(authH, bearerPrefix)
	} else {
		authCookie, err := c.Cookie("_forge_auth")
		if err == nil {
			token = authCookie
		}
	}
	return strings.TrimSpace(token)
}
