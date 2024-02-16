package main

import (
	"log"
	"os"

	"{{ cookiecutter.module_prefix }}{{ cookiecutter.slug }}/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.ServerCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
