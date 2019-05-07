package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkgems/npm-package-registry/registry"
	"github.com/gin-gonic/gin"
)

// PackageRoutes provides routes for interactions with packages
func PackageRoutes(router *gin.Engine) {
	router.PUT("/:name", handlerWrapper(publishPackageHandler))
	router.GET("/:name", handlerWrapper(getPackageHandler))
	router.GET("/:name/-/:version", handlerWrapper(getPackageTarballHandler))
}

func publishPackageHandler(core *registry.Core, c *gin.Context) (int, error) {
	var user = registry.User{
		Username: "anonymous",
		Password: "anonymous",
	}

	data, err := c.GetRawData()
	if err != nil {
		return http.StatusBadRequest, err
	}

	err = core.PublishPackage(string(data), user, c.Request.URL.Scheme+c.Request.Host)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return 200, nil
}

func getPackageHandler(core *registry.Core, c *gin.Context) (int, error) {
	var name = c.Param("name")

	var pkg, err = core.GetPackage(name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if pkg == "" {
		return http.StatusNotFound, fmt.Errorf("package not found")
	}

	c.Data(http.StatusOK, "application/json", []byte(pkg))
	return 0, nil
}

func getPackageTarballHandler(core *registry.Core, c *gin.Context) (int, error) {
	var name = c.Param("name")
	var versionWithExt = c.Param("version")
	var version = strings.Replace(versionWithExt, ".tgz", "", 1)

	var tarball, err = core.GetPackageTarball(name, version)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if len(tarball) == 0 {
		return http.StatusNotFound, fmt.Errorf("package not found")

	}

	c.Data(http.StatusOK, "application/octet-stream", tarball)
	return 0, nil
}
