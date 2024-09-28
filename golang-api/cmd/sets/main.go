package main

import (
	"fmt"
	"time"

	"github.com/JorgeSaicoski/golang-volley-live-score/internal/api/handlers"
	"github.com/JorgeSaicoski/golang-volley-live-score/internal/services/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		fmt.Println("DB not Connected:", err)
		return
	} else {
		fmt.Println("DB Connected")
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ws", handlers.HandleWebSocket)
	r.GET("/ws/finish", handlers.FinishSet)

	r.Run(":8081")
}
