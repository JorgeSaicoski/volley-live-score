package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JorgeSaicoski/golang-volley-live-score/internal/services/database"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Operation string `json:"operation"`
	Team      string `json:"team"`
}

func HandleWebSocker(c *gin.Context) {
	role := c.DefaultQuery("role", "watch")
	matchIDParam := c.Query("matchID")
	matchID, err := strconv.Atoi(matchIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid matchID"})
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		fmt.Println("Failed to accept WebSocket connection:", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	if role == "interact" {
		handleInteractiveConnection(conn, uint(matchID))
	} else {
		handleWatchingConnection(conn)
	}
}

func handleWatchingConnection(conn *websocket.Conn) {
	ctx := context.Background()

	for {
		time.Sleep(2 * time.Second)
		err := wsjson.Write(ctx, conn, "match updated")
		if err != nil {
			fmt.Println("Error sending update to watcher:", err)
			return
		}
	}
}

func handleInteractiveConnection(conn *websocket.Conn, matchID uint) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for {
		var message Message
		err := wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Println("Read error or timeout:", err)
			return
		}
		switch message.Operation {
		case "add":
			fmt.Printf("Adding point to team %s\n", message.Team)
			updateScoreInSet(matchID, message.Team, "add")
		default:
			fmt.Printf("Removing point from team %s\n", message.Team)
			updateScoreInSet(matchID, message.Team, "remove")
		}

		var set database.Set
		if err := database.DB.Where("match_id = ? AND finished = ?", matchID, false).First(&set).Error; err == nil {
			scoreUpdate := gin.H{
				"team_a_score": set.ScoreTeamA,
				"team_b_score": set.ScoreTeamB,
			}
			wsjson.Write(ctx, conn, scoreUpdate)
		}

		fmt.Println("Received interactive message:", message)
	}
}

func updateScoreInSet(matchID uint, team string, operation string) error {
	var set database.Set

	if err := database.DB.Where("match_id = ? AND finished = ?", matchID, false).First(&set).Error; err != nil {
		return fmt.Errorf("failed to find active set: %w", err)
	}

	if team == "A" {
		if operation == "add" {
			set.ScoreTeamA++
		} else if set.ScoreTeamA > 0 {
			set.ScoreTeamA--
		}
	} else {
		if operation == "add" {
			set.ScoreTeamB++
		} else if set.ScoreTeamB > 0 {
			set.ScoreTeamB--
		}
	}

	if err := database.DB.Save(&set).Error; err != nil {
		return fmt.Errorf("failed to update set score: %w", err)
	}

	return nil
}

func FinishSet(c *gin.Context) {
	var set database.Set

	if err := database.DB.Where("finished = ?", false).First(&set).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching live set"})
		return
	}

	if set.Finished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Set already finished"})
		return
	}

	set.Finished = true
	if err := database.DB.Save(&set).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish the set"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Set finished", "set": set})
}
