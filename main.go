package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"

	"riggs/client"
	"riggs/serve"
)

// drone version number
var version string

func main() {
	app := cli.NewApp()
	app.Name = "riggs"
	app.Version = version
	app.Usage = "command line utility"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{}
	app.Commands = []cli.Command{
		serve.Command,
		client.Command,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
