package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stephenafamo/bob/gen"
	helpers "github.com/stephenafamo/bob/gen/bobgen-helpers"
	"github.com/stephenafamo/bob/gen/bobgen-psql/driver"
	"github.com/urfave/cli/v2"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	app := &cli.App{
		Name:      "bobgen-psql",
		Usage:     "Generate models and factories from your PostgreSQL database",
		UsageText: "bobgen-psql [-c FILE]",
		Version:   helpers.Version(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   helpers.DefaultConfigPath,
				Usage:   "Load configuration from `FILE`",
			},
		},
		Action: run,
	}

	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	configFile := c.String("config")

	config, driverConfig, err := helpers.GetConfig[driver.Config](configFile, "psql")
	if err != nil {
		return err
	}

	d := driver.New(driverConfig)

	cmdState := &gen.State[any]{
		Config:  config,
		Dialect: "psql",
	}

	return cmdState.Run(c.Context, d)
}
