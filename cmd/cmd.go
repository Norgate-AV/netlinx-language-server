package cmd

import "github.com/Norgate-AV/netlinx-language-server/internal/settings"

type Application struct {
	Serve   Serve
	options *settings.Options
	Verbose bool
}
