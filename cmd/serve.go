package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/internal/logger"
	"github.com/Norgate-AV/netlinx-language-server/internal/server"
	"github.com/Norgate-AV/netlinx-language-server/internal/transport"
	"github.com/Norgate-AV/netlinx-language-server/internal/workspace"
	"github.com/Norgate-AV/netlinx-language-server/parser"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func serve(ctx context.Context, c *cli.Command) error {
	if c.Bool("version") {
		fmt.Println(c.Version)
		return nil
	}

	log := setupLogging(c)

	ts, err := parser.NewTreeSitter()
	if err != nil {
		log.Error("Failed to create parser", logrus.Fields{
			"error": err.Error(),
		})

		return err
	}

	defer ts.Close()

	server, err := setupServer(log, ts)
	if err != nil {
		return err
	}

	defer func() {
		log.LogServerEvent("Shutting down server...")
		server.Stop()
	}()

	return startTransport(ctx, c, log, server)
}

func setupLogging(c *cli.Command) logger.Logger {
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

	return log
}

func setupServer(log logger.Logger, ts *parser.TreeSitter) (*server.Server, error) {
	state := workspace.NewState(&workspace.NewStateOptions{
		TreeSitter: ts,
		Logger:     log,
	})

	srv := server.NewServer(log, state)

	return srv, nil
}

func startTransport(_ context.Context, c *cli.Command, log logger.Logger, server *server.Server) error {
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

	defer func() {
		if err := t.Close(); err != nil {
			log.Error("Failed to close transport", logrus.Fields{
				"error": err.Error(),
				"type":  transportType,
			})
		}
	}()

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
