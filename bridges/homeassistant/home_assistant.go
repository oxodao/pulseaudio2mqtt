package bridges_homeassistant

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/mqtt2pulseaudio/config"
	"github.com/oxodao/mqtt2pulseaudio/utils"
)

type Availability struct {
	Topic string `json:"topic"`
}

type DiscoveryPacket struct {
	Name         string       `json:"name"`
	StateTopic   string       `json:"state_topic"`
	Availability Availability `json:"availability"`
}

func publishDiscovery(c mqtt.Client, prefix string) {
	data, _ := json.Marshal(DiscoveryPacket{
		Name:       config.GET.Broker.ClientID,
		StateTopic: utils.Topic("state"),
		Availability: Availability{
			Topic: utils.Topic("status"),
		},
	})

	t := c.Publish(prefix+"/sensor/"+config.GET.Broker.ClientID+"/config", 2, true, data)

	t.Wait()
}
