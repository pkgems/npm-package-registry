package adapter

import (
	"context"

	"encoding/base64"

	"github.com/tidwall/gjson"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StorageMongoDB provides implementation of Storage interface
// for storing data in memory
type StorageMongoDB struct {
	collection *mongo.Collection
}

type OptionsStorageMongoDB struct {
	Collection *mongo.Collection
}

// NewStorageMongoDB is used to initialize new StorageMongoDB
func NewStorageMongoDB(opts OptionsStorageMongoDB) *StorageMongoDB {
	// create index on name to improve lookup performance
	opts.Collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"name": 1},
		Options: &options.IndexOptions{
			Unique:     &[]bool{true}[0],
			Background: &[]bool{true}[0], // wtf, mongodb?
		},
	})

	return &StorageMongoDB{
		collection: opts.Collection,
	}
}

func (storage StorageMongoDB) WriteTarball(name, version string, data []byte) error {
	_, err := storage.collection.UpdateOne(context.TODO(), bson.M{"name": name + "-" + version}, bson.M{"$set": bson.M{"tarball": data}}, &options.UpdateOptions{
		Upsert: &[]bool{true}[0],
	})

	return err
}

func (storage StorageMongoDB) ReadTarball(name, version string) ([]byte, error) {
	res := storage.collection.FindOne(context.TODO(), bson.M{"name": name + "-" + version})

	err := res.Err()
	if err != nil {
		return nil, err
	}

	data, err := res.DecodeBytes()
	if err != nil {
		// not found
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}

		// error
		return nil, err
	}

	doc := data.String()
	tarball := gjson.Get(doc, "tarball.$binary.base64").String()

	return base64.StdEncoding.DecodeString(tarball)
}
