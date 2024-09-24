package main

import (
	"fmt"

	"github.com/JorgeSaicoski/golang-volley-live-score/internal/api/routes"
	"github.com/JorgeSaicoski/golang-volley-live-score/internal/services/database"
	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
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
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{"Content-Type"},
	}))

	routes.SetupRoutes(r)

	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start the server:", err)
	}

}
