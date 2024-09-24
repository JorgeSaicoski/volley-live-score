package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/JorgeSaicoski/golang-volley-live-score/internal/services/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMatches(c *gin.Context) {
	var matches []database.Match
	var count int64
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	offset := (page - 1) * pageSize

	if err := database.DB.Model(&database.Match{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := database.DB.Preload("Sets").Offset(offset).Limit(pageSize).Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"count":   count,
	})
}

func GetMatchLive(c *gin.Context) {
	var match database.Match

	if err := database.DB.Preload("Sets").Where("is_live = ?", true).First(&match).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "No live match found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching live match"})
		return
	}

	c.JSON(http.StatusOK, match)
}

func GetMatchByID(c *gin.Context) {
	var match database.Match

	matchId := c.Param("id")

	if err := database.DB.Preload("Sets").First(&match, matchId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

func CreateMatch(c *gin.Context) {
	var match database.Match

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if database.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not established"})
		return
	}

	if err := database.DB.Create(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

func UpdateMatch(c *gin.Context) {
	var match database.Match

	matchId := c.Param("id")

	if err := database.DB.Preload("Sets").First(&match, matchId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)

}

func ToggleMatchLive(c *gin.Context) {
	var match database.Match
	matchID := c.Param("id")

	if err := database.DB.First(&match, matchID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	if match.IsLive {
		match.IsLive = false
		if err := database.DB.Save(&match).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update match status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Match live status updated", "match": match})
		return
	}

	var previousLiveMatch database.Match
	if err := database.DB.Where("is_live = ?", true).First(&previousLiveMatch).Error; err == nil {
		previousLiveMatch.IsLive = false
		if err := database.DB.Save(&previousLiveMatch).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change the currently live match status"})
			return
		}
	}

	match.IsLive = true
	if err := database.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update match status"})
		return
	}

	if len(match.Sets) == 0 {
		newSet := database.Set{
			MatchID:    match.ID,
			ScoreTeamA: 0,
			ScoreTeamB: 0,
			Finished:   false,
			Win:        false,
		}
		if err := database.DB.Create(&newSet).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create initial set"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Match live status updated", "match": match})
}
