package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/jtdevlin/album-api/internal/controller"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	router.GET("games", controller.GetGames)
	router.GET("games/:id", controller.GetGameById)
	router.PATCH("games/:id", controller.StartGame)
	router.POST("games", controller.CreateGame)
	router.PATCH("games/:id/users", controller.AddUserToGame)
	router.PATCH("games/:id/cards/:cardValue", controller.SelectCard)
	router.POST("/games/:id/prompt", controller.SetPromptForGameId)
	router.Run()
}
