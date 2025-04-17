package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/urfave/cli/v3"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func NewRootCommand() *cli.Command {
	app := &cli.Command{
		Name:      "netlinx-language-server",
		Usage:     "A language server for Netlinx",
		Version:   version,
		Copyright: "Copyright © 2025 Norgate AV",

		HideHelpCommand: true,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-file",
				Aliases: []string{"l"},
				Usage:   "Path to log file (overrides env: NETLINX_LSP_LOG_FILE)",
				Value:   "",
				Sources: cli.EnvVars("NETLINX_LSP_LOG_FILE"),
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
				Value:   "stdio",
				Sources: cli.EnvVars("NETLINX_LSP_TRANSPORT"),
				Validator: func(value string) error {
					validTransports := []string{"stdio", "pipe", "socket"}

					if slices.Contains(validTransports, value) {
						return nil
					}

					return fmt.Errorf("invalid transport type: %s (must be 'stdio', 'pipe', or 'socket')", value)
				},
				ValidateDefaults: true,
			},
			&cli.StringFlag{
				Name:    "pipe",
				Usage:   "Pipe name for transport type 'pipe'",
				Value:   "netlinx-language-server-pipe",
				Sources: cli.EnvVars("NETLINX_LSP_PIPE"),
			},
			&cli.StringFlag{
				Name:    "port",
				Usage:   "Port for transport type 'socket'",
				Value:   "8080",
				Sources: cli.EnvVars("NETLINX_LSP_PORT"),
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
