package handler

import (
	"net/http"

	"github.com/emeralt/npm-package-registry/registry"
	"github.com/gin-gonic/gin"
)

// Handler provides http interface for Registry
func Handler(core *registry.Core) http.Handler {
	router := gin.Default()

	PackageRoutes(router, core)

	return router
}
