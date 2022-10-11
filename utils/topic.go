package utils

import (
	"fmt"

	"github.com/oxodao/mqtt2pulseaudio/config"
)

func Topic(topic string) string {
	return fmt.Sprintf("pa2mqtt/%v/%v", config.GET.Broker.ClientID, topic)
}
