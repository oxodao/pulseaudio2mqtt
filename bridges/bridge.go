package bridges

import (
	"fmt"

	"github.com/oxodao/mqtt2pulseaudio/config"
)

type Bridge interface {
	Name() string
	FullName() string
	Register() error
}

var availableBridges = []Bridge{}

func AddBridge(b Bridge) {
	availableBridges = append(availableBridges, b)
}

func Load() error {
	fmt.Println("Loading bridges...")
	for _, b := range availableBridges {
		if config.GET.GetBridgeConfig(b.Name()).Enabled {
			if err := b.Register(); err != nil {
				fmt.Printf("\t- Failed to load bridge %v\n%v\n", b.FullName(), err)
			} else {
				fmt.Printf("\t- Bridge: %v loaded\n", b.FullName())
			}
		}
	}

	fmt.Println()

	return nil
}
