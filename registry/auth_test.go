package registry

import (
	"testing"

	"github.com/emeralt/npm-package-registry/adapter"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	var (
		core = Core{
			database: adapter.NewDatabaseMemory(),
			storage:  adapter.NewStorageMemory(),
		}

		user = User{
			Username: "tester@emeralt.org",
			Password: "tester",
		}

		token string
		err   error
	)

	t.Run("RegisterUser", func(t *testing.T) {
		err = core.RegisterUser(user)

		assert.Nil(t, err)
	})

	t.Run("LoginUser", func(t *testing.T) {
		token, err = core.LoginUser(user)

		assert.Nil(t, err)
		assert.NotNil(t, token)
	})

	t.Run("DecodeToken", func(t *testing.T) {
		username, err := core.DecodeToken(token)

		assert.Nil(t, err)
		assert.Equal(t, username, user.Username)
	})
}
