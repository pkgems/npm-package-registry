package registry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	for _, core := range MakeTestCores() {
		var (
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
}
