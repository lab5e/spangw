// Package main implements simple  gateway to Span
package main

import (
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
	"github.com/lab5e/spangw/pkg/emulator"
	"github.com/lab5e/spangw/pkg/gw"
)

func main() {
	var config gw.Parameters
	kong.Parse(&config)

	emulatorHandler := emulator.New()
	gwHandler, err := gw.Create(config, emulatorHandler)
	if err != nil {
		slog.Error("Error creating gateway", "error", err)
		os.Exit(2)
	}
	defer gwHandler.Stop()
	if err := gwHandler.Run(); err != nil {
		slog.Error("Could not run the gateway process", "error", err)
		os.Exit(2)
	}
}
