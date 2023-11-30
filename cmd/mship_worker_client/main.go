package main

import (
	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Usage:   "Path to config file",
			EnvVars: []string{"CONFIG"},
			Value:   "/etc/mship/config.yaml",
		},
	}
}
