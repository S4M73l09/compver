package main

import (
	"os"

	"github.com/S4M73l09/compver/internal/cli"
)

func main() {
	app := cli.New()
	os.Exit(app.Run(os.Args[1:]))
}
