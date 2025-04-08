package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/lsp"
	"github.com/Norgate-AV/netlinx-language-server/rpc"
)

func main() {
	logger := getLogger("netlinx-language-server.log")
	logger.Println("Starting Netlinx Language Server...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v\n", err)
			continue
		}

		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
	logger.Printf("Received method: %s\n", method)
	logger.Printf("Received content: %s\n", content)

	switch method {
	case "initialize":
		logger.Println("Handling initialize method")

		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v\n", err)
			return
		}

		logger.Printf("Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)
	}
}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[netlinx-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}
