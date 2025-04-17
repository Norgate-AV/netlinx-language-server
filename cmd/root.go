package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/Norgate-AV/netlinx-language-server/internal/config"

	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func NewRootCommand() *cli.Command {
	app := &cli.Command{
		Name:      config.AppName,
		Usage:     config.AppDescription,
		Version:   version,
		Copyright: config.Copyright,

		HideHelpCommand: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-file",
				Aliases: []string{"l"},
				Usage:   "Path to log file",
				Value:   "",
				Sources: cli.EnvVars(config.EnvLogFile),
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging",
				Value: false,
			},
			&cli.StringFlag{
				Name:    "transport",
				Aliases: []string{"t"},
				Usage:   "Transport type (stdio, pipe, socket)",
				Value:   config.DefaultTransportType,
				Sources: cli.EnvVars(config.EnvTransport),
				Validator: func(value string) error {
					if slices.Contains(config.ValidTransports, value) {
						return nil
					}

					return fmt.Errorf("invalid transport type: %s (must be one of: %v)",
						value, config.ValidTransports)
				},
				ValidateDefaults: true,
			},
			&cli.StringFlag{
				Name:    "pipe",
				Usage:   "Pipe name for transport type 'pipe'",
				Value:   config.DefaultPipeName,
				Sources: cli.EnvVars(config.EnvPipe),
			},
			&cli.StringFlag{
				Name:    "port",
				Usage:   "Port for transport type 'socket'",
				Value:   config.DefaultPort,
				Sources: cli.EnvVars(config.EnvPort),
				Validator: func(value string) error {
					if _, err := strconv.Atoi(value); err != nil {
						return fmt.Errorf("invalid port number: %s", value)
					}

					return nil
				},
				ValidateDefaults: true,
			},
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			return serve(ctx, cmd)
		},
	}

	if commit != "" && date != "" {
		app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	}

	return app
}

func Execute(ctx context.Context) error {
	return NewRootCommand().Run(ctx, os.Args)
}
