package main

import (
	"os"

	"github.com/alex-jienexa/pad_template/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
