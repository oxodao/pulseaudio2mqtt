package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/noisetorch/pulseaudio"
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
		Topic    string `yaml:"topic"`
	} `yaml:"broker"`

	PlayerctlPath string             `yaml:"-"`
	PAClient      *pulseaudio.Client `yaml:"-"`

	DeltaPercentage float32       `yaml:"delta_percentage"`
	DeltaDuration   time.Duration `yaml:"delta_duration"`
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

	client, err := pulseaudio.NewClient()
	if err != nil {
		return err
	}

	cfg.PAClient = client

	GET = cfg
	return nil
}
