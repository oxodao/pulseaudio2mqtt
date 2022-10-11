package messages

import (
	"fmt"

	"github.com/oxodao/mqtt2pulseaudio/services"
	"github.com/oxodao/mqtt2pulseaudio/utils"
)

type MuteMessage struct {
	Mode string `json:"mode"`
}

func (m MuteMessage) CheckArgs() error {
	if !utils.Contains([]string{"TOGGLE", "ON", "OFF"}, m.Mode) {
		return fmt.Errorf("unknown mode: %v", m.Mode)
	}

	return nil
}

func (m MuteMessage) Process() {
	var err error = nil

	switch m.Mode {
	case "TOGGLE":
		_, err = services.GET.PAClient.ToggleMute()

	default:
		err = services.GET.PAClient.SetMute(m.Mode == "ON")
	}

	if err != nil {
		fmt.Println("Failed to set mute: ", err)
	}
}
