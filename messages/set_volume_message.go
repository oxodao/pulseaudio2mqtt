package messages

import (
	"fmt"
	"strings"

	"github.com/oxodao/mqtt2pulseaudio/config"
	"github.com/oxodao/mqtt2pulseaudio/services"
	"github.com/oxodao/mqtt2pulseaudio/utils"
)

type SetVolumeMessage struct {
	Mode  string   `json:"mode"`
	Value *float32 `json:"value_percentage,omitempty"`
}

func (m *SetVolumeMessage) MovementName() string {
	// The movement name should have an optional "sink/source" name in it
	return "set_volume"
}

func (m SetVolumeMessage) CheckArgs() error {
	if strings.HasPrefix(m.Mode, "VOL_") {
		if m.Value == nil {
			return fmt.Errorf("%v message has no value", m.Mode)
		}
	} else if !utils.Contains([]string{"CT_VOL_UP", "CT_VOL_DOWN", "CT_STOP"}, m.Mode) {
		return fmt.Errorf("unknown mode for a set_volume message: %v", m.Mode)
	}

	return nil
}

func (m *SetVolumeMessage) continuousMessageProcess(multiplier float32) {
	if m.Mode == "CT_STOP" {
		services.GET.StopMovement(m.MovementName())
		return
	}

	err := services.GET.StartMovement(m.MovementName(), services.NewMovement(config.GET.ContinuousDeltas.Delay, func() {
		currVol, err := services.GET.PAClient.Volume()
		if err != nil {
			fmt.Println(err)
			services.GET.StopMovement(m.MovementName())
			return
		}

		services.GET.PAClient.SetVolume(currVol + (config.GET.ContinuousDeltas.DeltaPercentage/100)*multiplier)
	}))

	if err != nil {
		fmt.Println("Failed to start movement: ", err)
	}
}

func (m *SetVolumeMessage) simpleMessageProcess(multiplier float32) {
	currVol, err := services.GET.PAClient.Volume()
	if err != nil {
		fmt.Println("Failed to get PA volume: ", err)
		return
	}

	perc := (*m.Value / 100) * multiplier
	services.GET.PAClient.SetVolume(currVol + perc)
}

func (m SetVolumeMessage) Process() {
	var mult int32 = 1
	if strings.HasSuffix(m.Mode, "_DOWN") {
		mult = -1
	}

	if strings.HasPrefix(m.Mode, "CT_") {
		m.continuousMessageProcess(float32(mult))
	} else if utils.Contains([]string{"VOL_UP", "VOL_DOWN"}, m.Mode) {
		m.simpleMessageProcess(float32(mult))
	} else if m.Mode == "VOL_SET" {
		services.GET.PAClient.SetVolume(*m.Value / 100)
	}
}
