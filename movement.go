package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/oxodao/mqtt2pulseaudio/config"
)

var mtx = sync.Mutex{}
var currentMovement *Movement = nil

type Movement struct {
	StartedAt time.Time
	Action    string
}

func processMovement() bool {
	mtx.Lock()
	defer mtx.Unlock()

	if currentMovement == nil {
		return true
	}

	currVol, err := config.GET.PAClient.Volume()
	if err != nil {
		fmt.Println("FAIL / ", err)
		return true
	}

	diff := config.GET.DeltaPercentage
	if currentMovement.Action == ACTION_VOL_MINUS {
		diff *= -1
	}

	config.GET.PAClient.SetVolume(currVol + diff)
	return false
}

func NewMovement(action string) {
	currentMovement = &Movement{
		StartedAt: time.Now(),
		Action:    action,
	}

	go func() {
		for {
			shouldStop := processMovement()
			time.Sleep(config.GET.DeltaDuration * time.Millisecond)

			if shouldStop {
				break
			}
		}
	}()
}

func StopMovement() {
	if currentMovement == nil {
		return
	}

	mtx.Lock()
	defer mtx.Unlock()

	currentMovement = nil
}
