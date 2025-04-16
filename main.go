package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/Norgate-AV/netlinx-language-server/internal/transport"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "netlinx-language-server"
	app.Usage = "A language server for Netlinx"
	app.Version = version

	app.HideVersion = true
	app.HideHelpCommand = true

	if commit != "" && date != "" {
		app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "log-file",
			Aliases: []string{"l"},
			Usage:   "Path to log file (overrides env: NETLINX_LSP_LOG_FILE)",
			Value:   "",
			EnvVars: []string{"NETLINX_LSP_LOG_FILE"},
		},
		&cli.BoolFlag{
			Name:  "verbose",
			Usage: "Enable verbose logging",
			Value: false,
		},
		&cli.BoolFlag{
			Name:               "version",
			Aliases:            []string{"v"},
			Usage:              "Print version information",
			DisableDefaultText: true,
		},
		&cli.StringFlag{
			Name:    "transport",
			Aliases: []string{"t"},
			Usage:   "Transport type (stdio, pipe, socket)",
			Value:   "stdio",
			EnvVars: []string{"NETLINX_LSP_TRANSPORT"},
			// Validator: func(value string) error {
			// 	if value != "stdio" && value != "pipe" && value != "socket" {
			// 		return fmt.Errorf("invalid transport type: %s", value)
			// 	}
			// 	return nil
			// },
			// ValidateDefaults: true,
		},
		&cli.StringFlag{
			Name:    "pipe",
			Usage:   "Pipe name for transport type 'pipe'",
			Value:   "netlinx-language-server-pipe",
			EnvVars: []string{"NETLINX_LSP_PIPE"},
		},
		&cli.StringFlag{
			Name:    "port",
			Usage:   "Port for transport type 'socket'",
			Value:   "8080",
			EnvVars: []string{"NETLINX_LSP_PORT"},
		},
	}

	app.Action = func(c *cli.Context) error {
		return serve(c)
	}

	app.Commands = []*cli.Command{}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func serve(c *cli.Context) error {
	if c.Bool("version") {
		fmt.Println(c.App.Version)
		return nil
	}

	logFile := c.String("log-file")
	if logFile == "" {
		logFile = logger.GetDefaultLogPath()
	}

	if err := logger.EnsureLogDirectoryExists(logFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
	}

	log, err := logger.NewFileLogger(logFile)
	if err != nil {
		log = logger.NewStdLogger()
		log.Error("Failed to initialize file logger", logrus.Fields{
			"error":    err.Error(),
			"fallback": "stderr",
		})
	}

	if c.Bool("verbose") {
		if l := logger.GetLogrusLogger(log); l != nil {
			l.SetLevel(logrus.DebugLevel)
			log.Info("Verbose logging enabled", logrus.Fields{
				"level": "debug",
			})
		}
	}

	ts, err := parser.NewTreeSitter()
	if err != nil {
		log.Error("Failed to create parser", logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	defer ts.Close()

	state := analysis.NewState(&analysis.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	server := server.NewServer(log, state)
	defer func() {
		log.LogServerEvent("Shutting down server...")
		server.Stop()
	}()

	// var connOpt []jsonrpc2.ConnOpt
	// if trace {
	// 	connOpt = append(connOpt, jsonrpc2.LogMessages(log.New(logWriter, "", 0)))
	// }

	transportType := c.String("transport")
	log.Info("Using transport", logrus.Fields{
		"type": transportType,
	})

	t, err := transport.CreateTransport(transportType, &transport.Options{
		PipeName:   c.String("pipe"),
		SocketPort: c.String("port"),
		Logger:     log,
	})
	if err != nil {
		log.Error("Failed to create transport", logrus.Fields{
			"error": err.Error(),
			"type":  transportType,
		})

		return err
	}

	defer t.Close()

	log.LogServerEvent(fmt.Sprintf("Starting server with %s transport", transportType))

	disconnectChan, err := t.Start(context.Background(), server)
	if err != nil {
		log.Error("Failed to start transport", logrus.Fields{
			"error": err.Error(),
			"type":  transportType,
		})

		return err
	}

	log.LogServerEvent("Started Netlinx Language Server...")

	<-disconnectChan
	log.LogServerEvent("Connections closed")

	return nil
}
