package cli

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/emeralt/npm-package-registry/adapter"
	"github.com/emeralt/npm-package-registry/handler"
	"github.com/emeralt/npm-package-registry/registry"
)

// start a server
var rootCmd = &cobra.Command{
	Use:   "npr",
	Short: "npr is a tiny npm package registry",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	core, err := registry.NewCore(registry.CoreConfig{
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
