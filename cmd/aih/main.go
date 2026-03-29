package main

import (
	"os"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/cli/aih"
)

func main() {
	app := aih.New()
	os.Exit(app.Run(os.Args[1:]))
}
