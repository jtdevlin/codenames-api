package model

import (
	"encoding/json"
	"fmt"
)

type CardType int

const (
	Unassigned CardType = iota
	Assassin
	Blue
	Red
	Civilian
)

var cardName = map[CardType]string{
	Unassigned: "unassigned",
	Assassin:   "assassin",
	Blue:       "blue",
	Red:        "red",
	Civilian:   "civilian",
}

func (ct CardType) String() string {
	return cardName[ct]
}

func (ct CardType) OtherTeam() CardType {
	if ct == Blue {
		return Red
	} else if ct == Red {
		return Blue
	}
	return Unassigned
}
func (ct *CardType) FromString(cardType string) CardType {
	return map[string]CardType{
		"unassigned": Unassigned,
		"assassin":   Assassin,
		"blue":       Blue,
		"red":        Red,
		"civilian":   Civilian,
	}[cardType]
}

func (ct CardType) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.String())
}

func (ct *CardType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("error unmarshalling json: " + err.Error())
		return err
	}
	fmt.Println("String value: " + s)
	*ct = ct.FromString(s)
	return nil
}

type Card struct {
	Value    string   `json:"value"`
	Type     CardType `json:"type"`
	Selected bool     `json:"selected"`
}

func NewCard(value string, cardType CardType) *Card {
	return &Card{
		Value: value,
		Type:  cardType,
	}
}
