package model

import (
	"time"
)

type Game struct {
	ID               string           `json:"id"`
	Users            []*User          `json:"users"`
	CreatedTimestamp time.Time        `json:"createdTimestamp"`
	Cards            map[string]*Card `json:"cards"`
	State            GameState        `json:"state"`
	Turn             CardType         `json:"turn"`
	Winner           CardType         `json:"winner"`
}

type GameState string

const (
	Created   GameState = "Created"
	Ready     GameState = "Ready"
	Started   GameState = "Started"
	Completed GameState = "Completed"
)
