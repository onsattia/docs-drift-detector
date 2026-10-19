package main

import (
	"os"

	"github.com/onsattia/docs-drift-detector/internal/cli"
)

// version is overridden at build time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	os.Exit(cli.Execute(os.Args[1:], os.Stdout, os.Stderr, version))
}
