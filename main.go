package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Norgate-AV/netlinx-language-server/cmd"
)

func main() {
	if err := cmd.Execute(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
