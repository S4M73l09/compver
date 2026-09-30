package main

import (
	"fmt"
	"os"

	"github.com/S4M73l09/compver/internal/engine"
)

func main() {
	root := "."

	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	analyzer := engine.New()
	result, err := analyzer.Analyze(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Ruta analizada: %s\n", result.Path)
	fmt.Printf("Dependencias encontradas: %d\n", len(result.Dependencies))
}
