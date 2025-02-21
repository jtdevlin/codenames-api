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
	service.Games[newGame.ID] = &newGame
	context.IndentedJSON(http.StatusCreated, newGame)
}

func AddUserToGame(context *gin.Context) {
	gameID := context.Param("id")
	var user model.User

	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Add user request is malformed")
		return
	}

	game, ok := service.Games[gameID]
	if !ok {
		context.IndentedJSON(http.StatusPreconditionFailed, fmt.Sprintf("Game ID does not exist: %s", gameID))
		return
	}
	game.Users = append(game.Users, user)
	context.IndentedJSON(http.StatusOK, game)
}
