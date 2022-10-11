package services

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/AlexanderGrom/go-event"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/noisetorch/pulseaudio"
	"github.com/oxodao/mqtt2pulseaudio/bridges"
	"github.com/oxodao/mqtt2pulseaudio/config"
)

var GET *Provider = nil

/**
 *	Events available:
 *		- start: func()
 */

type Provider struct {
	MqttClient mqtt.Client
	Movements  map[string]*Movement
	PAClient   *pulseaudio.Client

	EventManager event.Dispatcher
	mvtMutex     sync.Mutex
}

func (p *Provider) getMovement(name string) *Movement {
	if m, ok := p.Movements[name]; ok {
		return m
	}

	return nil
}

func (p *Provider) StartMovement(name string, mvt *Movement) error {
	p.mvtMutex.Lock()
	if currMvt := p.getMovement(name); currMvt != nil && currMvt.Running {
		p.mvtMutex.Unlock()
		return errors.New("a movement with this name is already running")
	}

	p.Movements[name] = mvt.Run()
	p.mvtMutex.Unlock()

	return nil
}

func (p *Provider) StopMovement(name string) error {
	p.mvtMutex.Lock()

	mvt := p.getMovement(name)
	if mvt == nil || !mvt.Running {
		p.mvtMutex.Unlock()
		return errors.New("tried to stop a movement that was not running")
	}

	mvt.Stop()
	p.mvtMutex.Unlock()

	return nil
}

func Load() error {
	client := mqtt.NewClient(config.GET.GetClientConfig())

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	// Ugly hack to have a correct output
	// It seems the OnConnected handler event is slower than subscribing to everything
	time.Sleep(1 * time.Second)

	paClient, err := pulseaudio.NewClient()
	if err != nil {
		return err
	}

	GET = &Provider{
		MqttClient:   client,
		Movements:    map[string]*Movement{},
		mvtMutex:     sync.Mutex{},
		PAClient:     paClient,
		EventManager: event.New(),
	}

	err = bridges.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
