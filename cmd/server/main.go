package main

import (
	"github.com/muniere/glean/internal/app/server/cli"
	"github.com/muniere/glean/internal/pkg/sys"
)

func main() {
	cmd := cli.NewCommand()
	err := cmd.Execute()
	sys.CheckError(err)
}
