package registry

import (
	"testing"

	"github.com/emeralt/npm-registry/adapter"
	"github.com/stretchr/testify/assert"
)

var (
	core = Core{
		database: adapter.NewDatabaseMemory{},
	}

	user = User{
		Username: "tester@emeralt.org",
		Password: "tester",
	}

	token string
	err   error
)

func TestCore_RegisterUser(t *testing.T) {
	err := core.RegisterUser(user)

	assert.Nil(t, err)
}

func TestCore_LoginUser(t *testing.T) {
	token, err = core.LoginUser(user)

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestCore_DecodeToken(t *testing.T) {
	username, err := core.DecodeToken(token)

	assert.Nil(t, err)
	assert.Equal(t, username, user.Username)
}
