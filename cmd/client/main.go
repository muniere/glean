package main

import (
	"github.com/muniere/glean/internal/app/client"
	"github.com/muniere/glean/internal/pkg/sys"
)

func main() {
	cmd := client.NewCommand()
	err := cmd.Execute()
	sys.CheckError(err)
}
