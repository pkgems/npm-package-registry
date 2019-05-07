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
