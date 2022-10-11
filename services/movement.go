package services

import (
	"sync"
	"time"
)

type Movement struct {
	StartedAt time.Time
	LoopDelta time.Duration
	Running   bool
	Action    func()

	mtx    sync.Mutex
	onStop func(duration time.Duration)
}

func (m *Movement) Run() *Movement {
	m.Running = true

	go func() {
		for {
			m.mtx.Lock()
			m.Action()
			m.mtx.Unlock()
			time.Sleep(m.LoopDelta * time.Millisecond)

			if !m.Running {
				if m.onStop != nil {
					m.onStop(time.Since(m.StartedAt))
				}

				return
			}
		}
	}()

	return m
}

func (m *Movement) Stop() *Movement {
	m.Running = false

	return m
}

func (m *Movement) OnStop(onStop func(time.Duration)) *Movement {
	m.onStop = onStop

	return m
}

func NewMovement(loopDelta time.Duration, action func()) *Movement {
	return &Movement{
		StartedAt: time.Now(),
		LoopDelta: loopDelta,
		Running:   false,
		Action:    action,
		mtx:       sync.Mutex{},
	}
}
