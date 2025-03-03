package model

import (
	"time"
)

type Prompt struct {
	Clue   string `json:"clue"`
	Number int    `json:"number"`
}

type Game struct {
	ID                 string           `json:"id"`
	Users              []*User          `json:"users"`
	CreatedTimestamp   time.Time        `json:"createdTimestamp"`
	Cards              map[string]*Card `json:"cards"`
	State              GameState        `json:"state"`
	Turn               CardType         `json:"turn"`
	Winner             CardType         `json:"winner"`
	Prompt             *Prompt          `json:"prompt"`
	GuessesRemaining   int              `json:"guessesRemaining"`
	BlueCardsRemaining int              `json:"blueCardsRemaining"`
	RedCardsRemaining  int              `json:"redCardsRemaining"`
}

type GameState string

const (
	Created   GameState = "Created"
	Ready     GameState = "Ready"
	Started   GameState = "Started"
	Completed GameState = "Completed"
)
