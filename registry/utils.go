package registry

import "github.com/pkgems/npm-package-registry/adapter"

// MakeTestCores is internal function used to make a range of cores
// using all available database and storage adapters
func MakeTestCores() []Core {
	var cores []Core

	var databases = []adapter.Database{
		adapter.NewDatabaseMemory(),
	}

	var storages = []adapter.Storage{
		adapter.NewStorageMemory(),
	}

	for _, database := range databases {
		for _, storage := range storages {
			core := Core{
				database: database,
				storage:  storage,
			}

			cores = append(cores, core)
		}
	}

	return cores
}
