package main

const (
	ACTION_PRESS     = "toggle"
	ACTION_VOL_MINUS = "brightness_move_down"
	ACTION_VOL_PLUS  = "brightness_move_up"
	ACTION_VOL_STOP  = "brightness_stop"
)

type Message struct {
	Action          string `json:"action"`
	Rate            int    `json:"action_rate"`
	Battery         int    `json:"battery"`
	LinkQuality     int    `json:"linkquality"`
	UpdateAvailable bool   `json:"update_available"`
	Update          struct {
		State string `json:"state"`
	} `json:"update"`
}
