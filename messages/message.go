package messages

import (
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/mqtt2pulseaudio/services"
	"github.com/oxodao/mqtt2pulseaudio/utils"
)

type Message interface {
	Process()
	CheckArgs() error
}

func subscribe(channel string, receive func(mqtt.Client, mqtt.Message), desc string) error {
	channel = utils.Topic(channel)

	fmt.Printf("\t[%v] %v\n", channel, desc)
	token := services.GET.MqttClient.Subscribe(channel, 2, receive)
	token.Wait()

	return token.Error()
}

func Subscribe() error {
	fmt.Println("Subscribing to: ")
	subscribe("set_volume", Receive[SetVolumeMessage](), "Sets the PulseAudio volume")
	subscribe("mute", Receive[MuteMessage](), "Toggle the mute")

	fmt.Println()

	return nil
}

func Receive[T Message]() func(mqtt.Client, mqtt.Message) {
	return func(c mqtt.Client, m mqtt.Message) {
		data := new(T)
		err := json.Unmarshal(m.Payload(), &data)
		if err != nil {
			fmt.Println("Failed to parse message: ", err)
			return
		}

		if err = (*data).CheckArgs(); err != nil {
			fmt.Println("Failed to process message: missing arg", err)
			return
		}

		(*data).Process()
	}
}
