package bridges_homeassistant

import (
	"fmt"
	"time"

	"github.com/oxodao/mqtt2pulseaudio/config"
	"github.com/oxodao/mqtt2pulseaudio/services"
	"github.com/oxodao/mqtt2pulseaudio/utils"
)

type HomeAssistantBridge struct {
	LastSentValue float32
}

func (b HomeAssistantBridge) GetConfig() map[string]interface{} {
	return config.GET.GetBridgeConfig(b.Name()).Settings
}

func (b HomeAssistantBridge) Name() string {
	return "home_assistant"
}

func (b HomeAssistantBridge) FullName() string {
	return "Home Assistant"
}

func (b HomeAssistantBridge) GetPrefix() string {
	if val, ok := b.GetConfig()["autodiscovery_prefix"]; ok {
		return fmt.Sprint(val)
	}

	return "homeassistant"
}

func (b HomeAssistantBridge) Register() error {
	b.LastSentValue = -1

	services.GET.EventManager.On("start", func() error {
		publishDiscovery(services.GET.MqttClient, b.GetPrefix())
		return nil
	})

	go func() {
		for {
			val, err := services.GET.PAClient.Volume()
			if err != nil {
				fmt.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}

			if !utils.FloatEquals(val, b.LastSentValue, 0.01) {
				t := services.GET.MqttClient.Publish(utils.Topic("state"), 2, false, "10")
				t.Wait()

				if t.Error() != nil {
					fmt.Println(t.Error())
				}
			}

			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}
