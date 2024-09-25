package routes

import (
	"github.com/JorgeSaicoski/golang-volley-live-score/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/matches", handlers.GetMatches)
	r.GET("/matches/live", handlers.GetMatchLive)
	r.GET("/matches/:id", handlers.GetMatchByID)
	r.PATCH("/matches/:id/live", handlers.ToggleMatchLive)
	r.POST("/matches", handlers.CreateMatch)
	r.POST("/matches/:id", handlers.CreateSet)
	r.PUT("/matches/:id", handlers.UpdateMatch)
}
