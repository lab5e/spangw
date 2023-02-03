// Package main implements simple  gateway to Span
package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lab5e/spangw/pkg/emulator"
	"github.com/lab5e/spangw/pkg/gw"
	"github.com/lab5e/spangw/pkg/lg"
)

func main() {
	var config gw.Parameters
	kong.Parse(&config)

	emulatorHandler := emulator.New()
	gwHandler, err := gw.Create(config, emulatorHandler)
	if err != nil {
		lg.Error("Error creating gateway: %v", err)
		os.Exit(2)
	}
	defer gwHandler.Stop()
	if err := gwHandler.Run(); err != nil {
		lg.Error("Could not run the gateway process: %v", err)
		os.Exit(2)
	}
}
