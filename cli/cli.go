//
// Package cli provides command line interface for npm-package-registry (NPR).
// See https://github.com/pkgems/npm-package-registry/ for more information about registry.
//
package cli

import (
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "[NPR] ", 0)

// Run starts the CLI
func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
