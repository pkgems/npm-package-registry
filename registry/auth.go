package registry

import (
	"encoding/base64"
)

// User provides base user structure
type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterUser is used to create a new user
func (core Core) RegisterUser(user User) error {
	return nil
}

// LoginUser is used to log existing user in and return signed token
func (core Core) LoginUser(user User) (string, error) {
	token := core.SignToken(user)

	return token, nil
}

// SignToken is used to sign authentication token for user
func (core Core) SignToken(user User) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Username))
}

// DecodeToken to decode authentication token and retrieve owner
func (core Core) DecodeToken(token string) (string, error) {
	username, err := base64.StdEncoding.DecodeString(token)

	return string(username), err
}
