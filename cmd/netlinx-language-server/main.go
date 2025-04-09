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
			Name:    "verbose",
			Aliases: []string{"v"},
			Usage:   "Enable verbose logging",
			Value:   false,
		},
	}

	app.Action = func(c *cli.Context) error {
		log, err := logger.NewFileLogger("netlinx-language-server.log")
		if err != nil {
			log = logger.NewStdLogger()
			log.Printf("Failed to initialize file logger: %v, falling back to stderr", err)
		}

		if c.Bool("verbose") {
			if l := logger.GetLogrusLogger(log); l != nil {
				l.SetLevel(logrus.DebugLevel)
				log.Println("Verbose logging enabled")
			}
		}

		log.Println("Started Netlinx Language Server...")

		state := analysis.NewState()

		// Create JSONRPC handler
		server := server.NewServer(log, state)

		// Create and run JSONRPC connection using standard input/output
		<-jsonrpc2.NewConn(
			context.Background(),
			jsonrpc2.NewBufferedStream(&stdinStdout{}, jsonrpc2.VSCodeObjectCodec{}),
			server,
		).DisconnectNotify()

		return nil
	}

	app.Commands = []*cli.Command{
		{
			Name:  "version",
			Usage: "Print the version of the language server",
			Action: func(c *cli.Context) error {
				fmt.Printf("Version: %s\n", version)
				if commit != "" && date != "" {
					fmt.Printf("Commit: %s\n", commit)
					fmt.Printf("Date: %s\n", date)
				}

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
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
