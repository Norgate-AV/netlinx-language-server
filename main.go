package main

import (
	"bufio"
	"log"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/rpc"
)

func main() {
	logger := getLogger("netlinx-language-server.log")
	logger.Println("Starting Netlinx Language Server...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	for scanner.Scan() {
		msg := scanner.Bytes()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg []byte) {
	logger.Printf("Received message: %s\n", msg)
}

func getLogger(fileName string) *log.Logger {
	logFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[netlinx-language-server]", log.Ldate|log.Ltime|log.Lshortfile)
}
