package registry

import (
	"encoding/json"
	"testing"

	"github.com/emeralt/npm-package-registry/adapter"
	"github.com/stretchr/testify/assert"
)

func TestPackage(t *testing.T) {
	var (
		core = Core{
			database: adapter.NewDatabaseMemory(),
			storage:  adapter.NewStorageMemory(),
		}

		user = User{
			Username: "tester@emeralt.org",
			Password: "tester",
		}

		pkg            = `{"name": "test", "versions": {"1.0.0": {"name": "test", "version": "1.0.0"}}, "_attachments": {"test-1.0.0.tgz": {}}}`
		pkgName        = "test"
		pkgVersion     = "1.0.0"
		registryURL    = "test"
		decodedPackage Package
	)

	t.Run("PublishPackage", func(t *testing.T) {
		assert.Error(t, core.PublishPackage("", user, registryURL))
		assert.Error(t, core.PublishPackage("{}", user, registryURL))
		assert.Error(t, core.PublishPackage(`{"name": "test", "versions": {}}`, user, registryURL))
		assert.Error(t, core.PublishPackage(`{"versions": {"1.0.0": {}}}`, user, registryURL))
		assert.Error(t, core.PublishPackage(`{"name": "test", "versions": {"1.0.0": {"name": "test", "version": "1.0.0"}}}`, user, registryURL))
		assert.Error(t, core.PublishPackage(`{"name": "test", "versions": {"1.0.0": {"name": "test", "version": "1.0.0"}}, "_attachments": {}}`, user, registryURL))
		assert.NoError(t, core.PublishPackage(pkg, user, ""))
	})

	t.Run("GetPackage", func(t *testing.T) {
		tPkg, err := core.GetPackage("")
		assert.NoError(t, err)
		assert.Equal(t, tPkg, "")

		tPkg, err = core.GetPackage("nonexistent")
		assert.NoError(t, err)
		assert.Equal(t, tPkg, "")

		tPkg, err = core.GetPackage(pkgName)
		assert.NoError(t, err)
		assert.NotEqual(t, tPkg, pkg)

		err = json.Unmarshal([]byte(tPkg), &decodedPackage)
		assert.NoError(t, err)
	})

	t.Run("GetTarball", func(t *testing.T) {
		tarball, err := core.GetPackageTarball(pkgName, pkgVersion)
		assert.NoError(t, err)
		assert.Equal(t, tarball, []byte(""))
	})
}
