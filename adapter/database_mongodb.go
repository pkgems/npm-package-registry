package adapter

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DatabaseMongoDB provides implementation of Database interface
// for storing data in memory
type DatabaseMongoDB struct {
	collection *mongo.Collection
}

// OptionsDatabaseMongoDB parameterizes DatabaseMongoDB
type OptionsDatabaseMongoDB struct {
	Collection *mongo.Collection
}

// NewDatabaseMongoDB is used to initialize new DatabaseMongoDB
func NewDatabaseMongoDB(opts OptionsDatabaseMongoDB) *DatabaseMongoDB {
	// create index on name to improve lookup performance
	opts.Collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"name": 1},
		Options: &options.IndexOptions{
			Unique:     &[]bool{true}[0],
			Background: &[]bool{true}[0], // wtf, mongodb?
		},
	})

	return &DatabaseMongoDB{
		collection: opts.Collection,
	}
}

func (database DatabaseMongoDB) SetPackage(name string, data string) error {
	var bdoc interface{}
	err := bson.UnmarshalExtJSON([]byte(data), true, &bdoc)
	if err != nil {
		return err
	}

	_, err = database.collection.UpdateOne(context.TODO(), bson.M{"name": name}, bson.M{"$set": bdoc}, &options.UpdateOptions{
		Upsert: &[]bool{true}[0], // wtf, mongodb?
	})

	return err
}

func (database DatabaseMongoDB) GetPackage(name string) (string, error) {
	res := database.collection.FindOne(context.TODO(), bson.M{"name": name})

	err := res.Err()
	if err != nil {
		return "", err
	}

	data, err := res.DecodeBytes()
	if err != nil {
		// not found
		if err.Error() == "mongo: no documents in result" {
			return "", nil
		}

		// error
		return "", err
	}

	return data.String(), nil
}
