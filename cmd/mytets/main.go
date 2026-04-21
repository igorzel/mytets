package main

import (
	"os"

	"github.com/igorzel/mytets/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
