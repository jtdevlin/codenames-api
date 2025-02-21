package main

import (
	"github.com/gin-gonic/gin"

	"github.com/jtdevlin/album-api/internal/controller"
)

func main() {
	router := gin.Default()
	router.GET("games", controller.GetGames)
	router.GET("games/:id", controller.GetGameById)
	router.POST("games", controller.CreateGame)
	router.PATCH("games/:id/users", controller.AddUserToGame)
	router.Run()
}
