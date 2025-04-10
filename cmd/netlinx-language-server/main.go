package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/analysis"
	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
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
			Usage:   "Path to log file",
			Value:   "netlinx-language-server.log",
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

	log.LogServerEvent("Started Netlinx Language Server...")

	state := analysis.NewState()

	server := server.NewServer(log, state)
	// defer func() {
	// log.Error("Error shutting down server", logrus.Fields{
	//     "error": err.Error(),
	// })
	// }()

	// var connOpt []jsonrpc2.ConnOpt
	// if trace {
	// 	connOpt = append(connOpt, jsonrpc2.LogMessages(log.New(logWriter, "", 0)))
	// }

	log.LogServerEvent("Reading from stdin, writing to stdout")

	<-jsonrpc2.NewConn(
		context.Background(),
		jsonrpc2.NewBufferedStream(&stdinStdout{}, jsonrpc2.VSCodeObjectCodec{}),
		server,
		// connOpt...,
	).DisconnectNotify()

	log.LogServerEvent("Connections closed")

	return nil
}

type stdinStdout struct{}

func (stdinStdout) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (stdinStdout) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (stdinStdout) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	return os.Stdout.Close()
}
