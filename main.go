package main

import (
	"log"
	"net/http"

	"github.com/emeralt/npm-package-registry/adapter"
	"github.com/emeralt/npm-package-registry/handler"
	"github.com/emeralt/npm-package-registry/registry"
)

func main() {
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
