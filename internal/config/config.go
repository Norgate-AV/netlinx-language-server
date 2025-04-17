package config

// Application information
const (
	AppName        = "netlinx-language-server"
	Copyright      = "Copyright © 2025 Norgate AV"
	AppDescription = "A language server for Netlinx"
)

// Transport options
const (
	TransportTypeStdio  = "stdio"
	TransportTypePipe   = "pipe"
	TransportTypeSocket = "socket"

	DefaultTransportType = TransportTypeStdio
	DefaultPipeName      = "netlinx-language-server-pipe"
	DefaultPort          = "8080"
)

// Environment variable names
const (
	EnvLogFile   = "NETLINX_LSP_LOG_FILE"
	EnvTransport = "NETLINX_LSP_TRANSPORT"
	EnvPipe      = "NETLINX_LSP_PIPE"
	EnvPort      = "NETLINX_LSP_PORT"
)

// Valid options
var ValidTransports = []string{TransportTypeStdio, TransportTypePipe, TransportTypeSocket}
