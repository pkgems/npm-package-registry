package handler

import (
	"net/http"
	"strings"

	"github.com/emeralt/npm-package-registry/registry"
	"github.com/gin-gonic/gin"
)

// PackageRoutes provides routes for interactions with packages
func PackageRoutes(router *gin.Engine, core *registry.Core) {
	router.PUT("/:name", publishPackageHandler(core))
	router.GET("/:name", getPackageHandler(core))
	router.GET("/:name/-/:version", getPackageTarballHandler(core))
}

func publishPackageHandler(core *registry.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user = registry.User{
			Username: "anonymous",
			Password: "anonymous",
		}

		data, err := c.GetRawData()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = core.PublishPackage(string(data), user, c.Request.URL.Scheme+c.Request.Host)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
}

func getPackageHandler(core *registry.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var name = c.Param("name")

		var pkg, err = core.GetPackage(name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if pkg == "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
		} else {
			c.Data(http.StatusOK, "application/json", []byte(pkg))
		}
	}
}

func getPackageTarballHandler(core *registry.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		var name = c.Param("name")
		var versionWithExt = c.Param("version")
		var version = strings.Replace(versionWithExt, ".tgz", "", 1)

		var tarball, err = core.GetPackageTarball(name, version)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if len(tarball) == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
		} else {
			c.Data(http.StatusOK, "application/octet-stream", tarball)
		}
	}
}
