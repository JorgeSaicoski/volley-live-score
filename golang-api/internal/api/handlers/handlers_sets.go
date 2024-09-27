package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
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

var scoreUpdateChan = make(chan gin.H)
var activeConnections = make(map[*websocket.Conn]bool)
var mu sync.Mutex

func HandleWebSocket(c *gin.Context) {
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
	handleNewConn(conn, uint(matchID))
	/*
		mu.Lock()
		activeConnections[conn] = true
		mu.Unlock()
	*/
	defer func() {
		mu.Lock()
		delete(activeConnections, conn)
		mu.Unlock()
	}()

	if role == "interact" {
		handleInteractiveConnection(conn, uint(matchID))
	} else {
		handleWatchingConnection()
	}
}

func handleWatchingConnection() {
	ctx := context.Background()

	for {
		select {
		case scoreUpdate := <-scoreUpdateChan:
			mu.Lock()
			for connection := range activeConnections {
				err := wsjson.Write(ctx, connection, scoreUpdate)
				if err != nil {
					fmt.Println("Error sending update to watcher:", err)
					connection.Close(websocket.StatusNormalClosure, "error")
					delete(activeConnections, connection)
				}
			}
			mu.Unlock()
		case <-time.After(2 * time.Minute):
			fmt.Println("No updates within 10 seconds, closing connection")
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
			mu.Lock()
			for connection := range activeConnections {
				err := wsjson.Write(ctx, connection, scoreUpdate)
				if err != nil {
					fmt.Println("Error sending update to watcher:", err)
					connection.Close(websocket.StatusNormalClosure, "error")
					delete(activeConnections, connection)
				}
			}
			mu.Unlock()
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
	scoreUpdate := gin.H{
		"team_a_score": set.ScoreTeamA,
		"team_b_score": set.ScoreTeamB,
	}
	scoreUpdateChan <- scoreUpdate

	return nil
}

func handleNewConn(conn *websocket.Conn, matchID uint) {
	mu.Lock()
	activeConnections[conn] = true
	ctx := context.Background()
	var set database.Set

	if err := database.DB.Where("match_id = ? AND finished = ?", matchID, false).First(&set).Error; err != nil {
		return
	}
	matchScore := gin.H{
		"team_a_score": set.ScoreTeamA,
		"team_b_score": set.ScoreTeamB,
	}

	err := wsjson.Write(ctx, conn, matchScore)
	if err != nil {
		fmt.Println("Error sending update to new connection:", err)
		conn.Close(websocket.StatusNormalClosure, "error")
		delete(activeConnections, conn)
	}
	mu.Unlock()

}
