package cli

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/pkgems/npm-package-registry/handler"
	"github.com/pkgems/npm-package-registry/registry"
)

var config struct {
	Silent        bool
	ListenAddress string
	Database      string
	Storage       string
}

func init() {
	rootCmd.Flags().BoolVar(&config.Silent, "silent", getEnvBool("SILENT", "0"), "Disable logging")
	rootCmd.Flags().StringVar(&config.ListenAddress, "listen", getEnvString("LISTEN_ADDRESS", "localhost:8080"), "Address to listen")
	rootCmd.Flags().StringVar(&config.Database, "database", getEnvString("DATABASE", "memory"), "Database to use. Available: memory, mongodb")
	rootCmd.Flags().StringVar(&config.Storage, "storage", getEnvString("STORAGE", "memory"), "Storage to use. Available: memory, mongodb")

	// database adapters
	// mongodb
	rootCmd.Flags().StringVar(&adaptersConfig.Database.MongoDB.URI, "database-mongodb-uri", getEnvString("DATABASE_MONGODB_URI", "mongodb://localhost:27017"), "MongoDB database connection uri")
	rootCmd.Flags().StringVar(&adaptersConfig.Database.MongoDB.Database, "database-mongodb-database", getEnvString("DATABASE_MONGODB_DATABASE", "pkgems"), "MongoDB database database")
	rootCmd.Flags().StringVar(&adaptersConfig.Database.MongoDB.Collection, "database-mongodb-collection", getEnvString("DATABASE_MONGODB_COLLECTION", "npm-package-registry"), "MongoDB database collection")

	// storage adapters
	// mongodb
	rootCmd.Flags().StringVar(&adaptersConfig.Storage.MongoDB.URI, "storage-mongodb-uri", getEnvString("STORAGE_MONGODB_URI", "mongodb://localhost:27017"), "MongoDB storage connection uri")
	rootCmd.Flags().StringVar(&adaptersConfig.Storage.MongoDB.Database, "storage-mongodb-database", getEnvString("STORAGE_MONGODB_DATABASE", "pkgems"), "MongoDB storage database")
	rootCmd.Flags().StringVar(&adaptersConfig.Storage.MongoDB.Collection, "storage-mongodb-collection", getEnvString("STORAGE_MONGODB_COLLECTION", "npm-package-registry"), "MongoDB storage collection")
}

// start a server
var rootCmd = &cobra.Command{
	Use:   "npr",
	Short: "npr is a tiny npm package registry",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logger.Printf("initializing database: %s", config.Database)
	database, err := getDatabase(config.Database)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Printf("initializing storage: %s", config.Storage)
	storage, err := getStorage(config.Storage)
	if err != nil {
		logger.Fatal(err)
	}

	core, err := registry.NewCore(registry.CoreConfig{
		Database: database,
		Storage:  storage,
	})
	if err != nil {
		log.Fatal(err)
	}

	server := http.Server{
		Handler: handler.Handler(core),
		Addr:    config.ListenAddress,
	}

	server.ListenAndServe()
}
