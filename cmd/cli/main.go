// cmd/cli/main.go
package main

import (
	"log"
	"os"

	"github.com/mohan2020coder/mSpace/internal/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
