package cli

import (
	"context"

	"github.com/pkgems/npm-package-registry/adapter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var adaptersConfig struct {
	Database struct {
		Memory  struct{}
		MongoDB struct {
			URI        string
			Database   string
			Collection string
		}
	}
	Storage struct {
		Memory  adapter.StorageMemory
		MongoDB struct {
			URI        string
			Database   string
			Collection string
		}
	}
}

func getDatabase(name string) (adapter.Database, error) {
	var database adapter.Database

	switch name {
	case "mongodb":
		// create connection
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(adaptersConfig.Database.MongoDB.URI))
		if err != nil {
			return nil, err
		}

		database = adapter.NewDatabaseMongoDB(adapter.OptionsDatabaseMongoDB{
			Collection: client.Database(adaptersConfig.Database.MongoDB.Database).Collection(adaptersConfig.Database.MongoDB.Collection),
		})
	default:
		database = adapter.NewDatabaseMemory()
	}

	return database, nil
}

func getStorage(name string) (adapter.Storage, error) {
	var storage adapter.Storage

	switch name {
	case "mongodb":
		// create connection
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(adaptersConfig.Storage.MongoDB.URI))
		if err != nil {
			return nil, err
		}

		storage = adapter.NewStorageMongoDB(adapter.OptionsStorageMongoDB{
			Collection: client.Database(adaptersConfig.Storage.MongoDB.Database).Collection(adaptersConfig.Storage.MongoDB.Collection),
		})
	default:
		storage = adapter.NewStorageMemory()
	}

	return storage, nil
}
