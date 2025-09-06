// cmd/cli/main.go
package main

import (
	"log"
	"os"

	"github.com/mohan2020coder/mSpace/internal/cli"
)

func main() {
	apiBase := "http://localhost:8080"

	if err := cli.Run(apiBase); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
