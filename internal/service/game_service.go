package service

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strings"

	"github.com/jtdevlin/album-api/internal/model"
)

var Games = map[string]*model.Game{}

func readCardValuesFromFile() ([]string, error) {
	filepath := "../internal/resource/words.txt"
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, errors.New("error reading words from file")
	}

	return strings.Split(string(content), "\n"), nil
}

func StartGame(gameId string) (model.Game, error) {
	game, ok := Games[gameId]
	if !ok {
		return model.Game{}, errors.New("No game ID found for ID: " + gameId)
	}
	if len(game.Users) < 4 {
		return model.Game{}, errors.New("game cannot be started with less than 4 players")
	}
	setTeamsForUsers(game)
	game.State = model.Started
	return *game, nil
}

func setTeamsForUsers(game *model.Game) {
	userSize := len(game.Users)
	redUsersRemaining := userSize / 2
	blueUsersRemaining := userSize - redUsersRemaining

	redUsers := make([]*model.User, 0)
	blueUsers := make([]*model.User, 0)

	for blueUsersRemaining > 0 {
		randomCardNumber := rand.Intn(userSize)
		currentUser := game.Users[randomCardNumber]
		if currentUser.Team == model.Unassigned {
			currentUser.Team = model.Blue
			blueUsers = append(blueUsers, currentUser)
			blueUsersRemaining--
		}
	}

	for redUsersRemaining > 0 {
		randomCardNumber := rand.Intn(userSize)
		currentUser := game.Users[randomCardNumber]
		if currentUser.Team == model.Unassigned {
			currentUser.Team = model.Red
			redUsers = append(redUsers, currentUser)
			redUsersRemaining--
		}
	}

	randomRedSpyMasterNumber := rand.Intn(len(redUsers))
	redUsers[randomRedSpyMasterNumber].IsSpyMaster = true
	randomBlueSpyMasterNumber := rand.Intn(len(blueUsers))
	blueUsers[randomBlueSpyMasterNumber].IsSpyMaster = true
}

func SetCardsForGame(game *model.Game) {
	allCards, err := readCardValuesFromFile()
	if err != nil {
		return
	}

	chosenCards := make(map[string]*model.Card)
	for i := 0; i < 25; i++ {
		randomNumber := rand.Intn(len(allCards))
		card := model.NewCard(allCards[randomNumber], model.Civilian)

		chosenCards[card.Value] = card
		//Remove card if it was already used
		allCards = slices.Delete(allCards, randomNumber, randomNumber+1)
	}
	firstTeam := chooseCardTypes(chosenCards)
	game.Turn = firstTeam
	game.Cards = chosenCards
}

func AddUserToGame(gameId string, user model.User) (model.Game, error) {

	game, ok := Games[gameId]
	if !ok {
		return model.Game{}, errors.New("game ID does not exist for ID: " + game.ID)
	}
	game.Users = append(game.Users, &user)
	if len(game.Users) >= 4 && game.State == model.Created {
		game.State = model.Ready
	}
	return *game, nil
}

func SelectedCard(gameId string, cardName string, user model.User) (model.Game, error) {
	game, ok := Games[gameId]
	if !ok {
		return model.Game{}, errors.New("No game ID found for ID: " + gameId)
	}
	card, ok := game.Cards[cardName]
	if !ok {
		return model.Game{}, errors.New("No card name found for name: " + cardName)
	}
	card.Selected = true
	return *game, nil
}

func chooseCardTypes(cards map[string]*model.Card) model.CardType {
	randomNumber := rand.Intn(2)
	blueCardsLeft := 8
	redCardsLeft := 8
	assassinCardsLeft := 1

	var firstTeam model.CardType
	if randomNumber == 1 {
		redCardsLeft++
		firstTeam = model.Red
	} else {
		blueCardsLeft++
		fmt.Println("Blue Goes First!")
		firstTeam = model.Blue
	}

	keys := make([]string, len(cards))
	i := 0
	for k := range cards {
		keys[i] = k
		i++
	}

	for redCardsLeft > 0 {
		randomCardNumber := rand.Intn(len(cards))
		currentKey := keys[randomCardNumber]
		currentCard := cards[currentKey]
		if currentCard.Type == model.Civilian {
			currentCard.Type = model.Red
			redCardsLeft--
		}
	}

	for blueCardsLeft > 0 {
		randomCardNumber := rand.Intn(len(cards))
		currentKey := keys[randomCardNumber]
		currentCard := cards[currentKey]
		if currentCard.Type == model.Civilian {
			currentCard.Type = model.Blue
			blueCardsLeft--
		}
	}

	for assassinCardsLeft > 0 {
		randomCardNumber := rand.Intn(len(cards))
		currentKey := keys[randomCardNumber]
		currentCard := cards[currentKey]
		if currentCard.Type == model.Civilian {
			currentCard.Type = model.Assassin
			assassinCardsLeft--
		}
	}

	return firstTeam
}
