package main

import (
	"fmt"
	"os"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/brokerdaemon"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("op-sa-broker %s\n", version.Version)
		return
	}

	cfg := config.Load()
	server := brokerdaemon.Server{Config: cfg}
	if err := server.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
