package model

import (
	"time"
)

type Game struct {
	ID               string    `json:"id"`
	Users            []User    `json:"users"`
	CreatedTimestamp time.Time `json:"createdTimestamp"`
	Cards            []string  `json:"cards"`
}
