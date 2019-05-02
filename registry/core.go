package registry

import (
	"github.com/emeralt/npm-package-registry/adapter"
)

// Core is used as the central manager of Registry activity. It is the primary point of
// interface for API handlers and is responsible for managing all the logic
type Core struct {
	database adapter.Database
	storage  adapter.Storage
	secret   string
}

// CoreConfig is used to parameterize the core
type CoreConfig struct {
	Database adapter.Database
	Storage  adapter.Storage
	Secret   string
}

// NewCore is used to construct a new core
func NewCore(conf CoreConfig) (*Core, error) {
	return &Core{
		database: conf.Database,
		storage:  conf.Storage,
		secret:   conf.Secret,
	}, nil
}
