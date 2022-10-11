package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gopkg.in/yaml.v2"
)

var GET *Config = nil

type Config struct {
	Broker struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		ClientID string `yaml:"client_id"`
	} `yaml:"broker"`

	Bridges map[string]BridgeConfig `json:"bridges"`

	ContinuousDeltas struct {
		Delay           time.Duration `yaml:"delay"`
		DeltaPercentage float32       `yaml:"delta_percentage"`
	} `yaml:"continuous_delta"`

	PlayerctlPath string `yaml:"-"`
}

type BridgeConfig struct {
	Enabled  bool
	Settings map[string]interface{}
}

func (c *Config) GetClientConfig() *mqtt.ClientOptions {
	mqttConfig := mqtt.
		NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%v:%v", c.Broker.Host, c.Broker.Port)).
		SetClientID(c.Broker.ClientID).
		SetOnConnectHandler(func(c mqtt.Client) {
			fmt.Println("Connected !")
		}).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			fmt.Printf("Connect lost: %v", err)
			os.Exit(1)
		})

	if len(c.Broker.Username) > 0 && len(c.Broker.Password) > 0 {
		mqttConfig = mqttConfig.
			SetUsername(c.Broker.Username).
			SetPassword(c.Broker.Password)
	}

	return mqttConfig
}

func (c *Config) GetBridgeConfig(name string) BridgeConfig {
	if val, ok := c.Bridges[name]; ok {
		return val
	}

	return BridgeConfig{
		Enabled:  false,
		Settings: map[string]interface{}{},
	}
}

func getConfigFile() string {
	filepaths := []string{"config.yaml", "/etc/mqtt2pulseaudio.yaml"}

	for _, fp := range filepaths {
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			continue
		}

		return fp
	}

	panic(fmt.Sprintf("FAILED TO FIND CONFIG FILE !! (Looked in %v)\n", strings.Join(filepaths, ", ")))
}

func Load() error {
	data, err := os.ReadFile(getConfigFile())
	if err != nil {
		return err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return err
	}

	pctl, err := exec.LookPath("playerctl")
	if err != nil {
		return err
	}

	cfg.PlayerctlPath = pctl

	GET = cfg
	return nil
}
