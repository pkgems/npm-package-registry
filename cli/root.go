package cli

import (
	"context"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkgems/npm-package-registry/adapter"
	"github.com/pkgems/npm-package-registry/handler"
	"github.com/pkgems/npm-package-registry/registry"
)

// start a server
var rootCmd = &cobra.Command{
	Use:   "npr",
	Short: "npr is a tiny npm package registry",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	core, err := registry.NewCore(registry.CoreConfig{
		// Database: adapter.NewDatabaseMongoDB(adapter.OptionsMongoDB{
		// 	Collection: client.Database("pkgems").Collection("npm-package-registry"),
		// }),
		Database: adapter.NewDatabaseMemory(),
		Storage:  adapter.NewStorageMemory(),
	})
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Handler: handler.Handler(core),
		Addr:    "localhost:8080",
	}

	server.ListenAndServe()
}
