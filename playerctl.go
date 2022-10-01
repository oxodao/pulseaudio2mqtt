package main

import (
	"os/exec"

	"github.com/oxodao/mqtt2pulseaudio/config"
)

func TogglePlay() error {
	cmd := exec.Command(config.GET.PlayerctlPath, "play-pause")
	return cmd.Run()
}
