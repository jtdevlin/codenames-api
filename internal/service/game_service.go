package service

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/jtdevlin/album-api/internal/model"
)

var Games = map[string]*model.Game{}

func readCardsFromFile() ([]string, error) {
	filepath := "../internal/resource/words.txt"
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, errors.New("error reading words from file")
	}

	return strings.Split(string(content), "\n"), nil
}

func SetCardsForGame(game *model.Game) {
	allCards, err := readCardsFromFile()
	if err != nil {
		return
	}

	chosenCards := make([]string, 0)
	for i := 0; i < 25; i++ {
		randomNumber := rand.Intn(len(allCards) - 1)
		chosenCards = append(chosenCards, allCards[randomNumber])
	}

	game.Cards = chosenCards
}
