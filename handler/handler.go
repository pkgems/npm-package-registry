package handler

import (
	"net/http"

	"github.com/pkgems/npm-package-registry/registry"
	"github.com/gin-gonic/gin"
)

// Handler provides http interface for Registry
func Handler(core *registry.Core) http.Handler {
	router := gin.Default()

	router.Use(coreMiddleware(core))

	PackageRoutes(router)

	return router
}

// DynamicHandler provides http interface for Registry
// allowing for dynamic core setup
func DynamicHandler(getCore func() *registry.Core) http.Handler {
	router := gin.Default()

	router.Use(dynamicCoreMiddleware(getCore))

	PackageRoutes(router)

	return router
}

func coreMiddleware(core *registry.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("core", core)

		c.Next()
	}
}

func dynamicCoreMiddleware(getCore func() *registry.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("core", getCore())

		c.Next()
	}
}

type HandlerFunc func(core *registry.Core, c *gin.Context) (int, error)

func handlerWrapper(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		core := c.MustGet("core").(*registry.Core)

		status, err := handler(core, c)
		if status != 0 {
			if err != nil {
				c.JSON(status, gin.H{
					"error": err.Error(),
				})
			} else {
				c.Status(status)
			}
		}
	}
}
