package model

type User struct {
	Name        string   `json:"name"`
	IsSpyMaster bool     `json:"isSpyMaster"`
	Team        CardType `json:"team"`
}
