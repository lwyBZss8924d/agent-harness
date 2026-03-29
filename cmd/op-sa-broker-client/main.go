package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/broker"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/version"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "version") {
		fmt.Printf("op-sa-broker-client %s\n", version.Version)
		return
	}

	action := "status"
	if len(os.Args) > 1 {
		action = os.Args[1]
	}

	cfg := config.Load()
	timeout := 5 * time.Second
	if action == "get-material" || action == "get_material" {
		timeout = 90 * time.Second
	}
	client := broker.Client{SocketPath: cfg.Broker.SocketPath, Timeout: timeout}

	switch action {
	case "status":
		response, err := client.Status()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(response)
	case "get-token", "get_token":
		response, err := client.GetToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if !response.OK {
			fmt.Fprintln(os.Stderr, response.Error)
			os.Exit(1)
		}
		fmt.Fprint(os.Stdout, response.Token)
	case "get-material", "get_material":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "usage: op-sa-broker-client get-material <reference>")
			os.Exit(2)
		}
		response, err := client.GetMaterial(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_ = json.NewEncoder(os.Stdout).Encode(response)
	default:
		fmt.Fprintf(os.Stderr, "unknown action: %s\n", action)
		os.Exit(2)
	}
}
