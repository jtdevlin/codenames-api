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

func GetGameByID(id string) (*model.Game, error) {
	game, exists := Games[id]
	if !exists {
		return nil, errors.New("game not found for ID: " + id)
	}
	return game, nil
}

func UpdateGame(game *model.Game) (*model.Game, error) {
	Games[game.ID] = game
	return game, nil
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
	blueCards, redCards := chooseCardTypes(chosenCards)
	game.BlueCardsRemaining = blueCards
	game.RedCardsRemaining = redCards
	if redCards > blueCards {
		game.Turn = model.Red
	} else {
		game.Turn = model.Blue
	}
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
	if game.Turn != user.Team {
		fmt.Println("It is not this user's turn! User team: " + user.Team.String())
		return model.Game{}, errors.New("It is not this user's turn")
	}
	if game.State == model.Completed {
		return model.Game{}, errors.New("The game is Over!")
	}
	card, ok := game.Cards[cardName]
	if !ok {
		return model.Game{}, errors.New("No card name found for name: " + cardName)
	}
	card.Selected = true
	//If the assassin card is selected, the other team wins
	if card.Type == model.Assassin {
		game.Winner = user.Team.OtherTeam()
		game.State = model.Completed
	} else if card.Type != user.Team {
		game.Turn = user.Team.OtherTeam()
		game.Prompt = nil
	} else {
		if user.Team == model.Blue {
			game.BlueCardsRemaining--
			if game.BlueCardsRemaining == 0 {
				game.Winner = model.Blue
				game.State = model.Completed
			}
		} else if user.Team == model.Red {
			game.RedCardsRemaining--
			if game.RedCardsRemaining == 0 {
				game.Winner = model.Red
				game.State = model.Completed
			}
		}
		game.GuessesRemaining--
		if game.GuessesRemaining == 0 {
			game.Turn = user.Team.OtherTeam()
			game.Prompt = nil
		}
	}

	return *game, nil
}

func chooseCardTypes(cards map[string]*model.Card) (int, int) {
	randomNumber := rand.Intn(2)
	assassinCardsLeft := 1

	totalRedCards, totalBlueCards := 8, 8
	if randomNumber == 1 {
		totalRedCards++
	} else {
		totalBlueCards++
	}

	keys := make([]string, len(cards))
	i := 0
	for k := range cards {
		keys[i] = k
		i++
	}

	blueCardsLeft := totalBlueCards
	redCardsLeft := totalRedCards
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

	return totalBlueCards, totalRedCards
}
