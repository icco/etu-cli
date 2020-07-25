package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "etu",
		Usage: "log a project to etu.natwelch.com",
		Commands: []*cli.Command{
			{
				Name:    "print",
				Aliases: []string{"p"},
				Usage:   "print recent entries",
				Action:  print,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a log",
				Action:  add,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "api_key",
				Usage:   "authorize your user",
				EnvVars: []string{"GQL_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "env",
				Usage:   "set which graphql server to talk to",
				Value:   "production",
				EnvVars: []string{"NAT_ENV"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func add(c *cli.Context) error {
	return nil
}

func print(c *cli.Context) error {
	return nil
}
