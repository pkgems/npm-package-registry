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

// OptionsMongoDB parameterizes DatabaseMongoDB
type OptionsMongoDB struct {
	Collection *mongo.Collection
}

// NewDatabaseMongoDB is used to initialize new DatabaseMongoDB
func NewDatabaseMongoDB(options OptionsMongoDB) *DatabaseMongoDB {
	return &DatabaseMongoDB{
		collection: options.Collection,
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
