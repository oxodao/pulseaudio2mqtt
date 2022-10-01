package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/oxodao/mqtt2pulseaudio/config"
)

func run() error {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := mqtt.NewClient(config.GET.GetClientConfig())

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token := client.Subscribe(config.GET.Broker.Topic, 2, func(c mqtt.Client, m mqtt.Message) {
		msg := Message{}
		err := json.Unmarshal(m.Payload(), &msg)
		if err != nil {
			fmt.Println("Failed to parse message: " + string(m.Payload()))
			return
		}

		switch msg.Action {
		case ACTION_PRESS:
			TogglePlay()
		case ACTION_VOL_MINUS:
			fallthrough
		case ACTION_VOL_PLUS:
			NewMovement(msg.Action)
		case ACTION_VOL_STOP:
			StopMovement()
		default:
			if len(msg.Action) > 0 {
				fmt.Println("Action pressed: " + msg.Action)
			}
		}
	})

	token.Wait()

	if token.Error() != nil {
		panic(token.Error())
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
