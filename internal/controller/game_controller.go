package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jtdevlin/album-api/internal/model"
	"github.com/jtdevlin/album-api/internal/service"
	"github.com/xyproto/randomstring"
)

// getGames responds with the list of all games as JSON.
func GetGames(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, service.Games)
}

func GetGameById(context *gin.Context) {
	id := context.Param("id")

	album, ok := service.Games[id]
	if ok {
		context.IndentedJSON(http.StatusOK, album)
	} else {
		context.IndentedJSON(http.StatusNotFound, fmt.Sprintf("Album not found for ID: %s", id))
	}
}

func CreateGame(context *gin.Context) {
	var newGame model.Game
	newGame.ID = randomstring.HumanFriendlyEnglishString(6)
	newGame.CreatedTimestamp = time.Now()

	service.SetCardsForGame(&newGame)
	newGame.State = model.Created
	service.Games[newGame.ID] = &newGame
	context.IndentedJSON(http.StatusCreated, newGame)
}

func StartGame(context *gin.Context) {
	gameId := context.Param("id")
	game, err := service.StartGame(gameId)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Unable to start game")
		return
	}

	context.IndentedJSON(http.StatusOK, game)
}

func AddUserToGame(context *gin.Context) {
	gameId := context.Param("id")
	var user model.User

	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Add user request is malformed")
		return
	}

	game, err := service.AddUserToGame(gameId, user)
	if err != nil {
		context.IndentedJSON(http.StatusPreconditionFailed, fmt.Sprintf("Game ID does not exist: %s", gameId))
		return
	}
	context.IndentedJSON(http.StatusOK, game)
}

func SelectCard(context *gin.Context) {
	gameID := context.Param("id")
	cardName := context.Param("cardValue")

	var user model.User

	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Provided user is malformed")
		return
	}

	game, err := service.SelectedCard(gameID, cardName, user)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, "No match for provided card and game ID")
	}
	context.IndentedJSON(http.StatusOK, game)
}

type SetPromptRequest struct {
	Clue   string `json:"clue" binding:"required"`
	Number int    `json:"number" binding:"required"`
}

func SetPromptForGameId(c *gin.Context) {
	gameID := c.Param("id")

	var req SetPromptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	game, err := service.GetGameByID(gameID)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	game.Prompt = &model.Prompt{
		Clue:   req.Clue,
		Number: req.Number,
	}
	game.GuessesRemaining = req.Number + 1
	updatedGame, err := service.UpdateGame(game)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game"})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedGame)
}
