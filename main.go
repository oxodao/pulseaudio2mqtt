package main

import (
	"fmt"
	"os"

	"github.com/oxodao/mqtt2pulseaudio/bridges"
	ha_bridge "github.com/oxodao/mqtt2pulseaudio/bridges/homeassistant"
	"github.com/oxodao/mqtt2pulseaudio/cmd"
	"github.com/oxodao/mqtt2pulseaudio/config"
)

func main() {
	bridges.AddBridge(ha_bridge.HomeAssistantBridge{})

	if err := config.Load(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.Execute()
}
