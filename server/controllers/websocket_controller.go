package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"forum/server/utils"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	ConnectedUsers = make(map[int]*Connection)
	mu             sync.Mutex
)

// Connection and user management
type Connection struct {
	Conn   *websocket.Conn
	UserID int
}

// HandleWebSocket manages a WebSocket connection for a user
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userID := r.Context().Value("user_id").(int)

	// Register the connection
	connection := &Connection{Conn: conn, UserID: userID}

	mu.Lock()
	ConnectedUsers[userID] = connection
	mu.Unlock()

	broadcastUserList()

	// Wait for disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// if the user disconnected
			break
		}
	}

	// Clean up on disconnection
	mu.Lock()
	delete(ConnectedUsers, userID)
	mu.Unlock()

	broadcastUserList()
}

// broadcastUserList sends the updated list of connected user IDs to all clients
func broadcastUserList() {
	mu.Lock()
	defer mu.Unlock()

	var userIDs []int
	for userID := range ConnectedUsers {
		userIDs = append(userIDs, userID)
	}

	// Broadcast to all connections
	for userID, connection := range ConnectedUsers {
		// Create a filtered list excluding the current user's ID
		filteredUserIDs := make([]int, 0, len(userIDs)-1)
		for _, id := range userIDs {
			if id != userID {
				filteredUserIDs = append(filteredUserIDs, id)
			}
		}

		message, err := json.Marshal(filteredUserIDs)
		if err != nil {
			log.Printf("Error marshalling user list for user %d: %v\n", userID, err)
			continue
		}

		err = connection.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// if there was an error it means that the user is disconnected
			connection.Conn.Close()
			delete(ConnectedUsers, userID)
		}
	}
}

func handleChat(userID int, conn *websocket.Conn) error {
	_, data, err := readMessage(userID, conn)
	log.Println(string(data))
	if err != nil {
		return fmt.Errorf("WebSocket read message failed: %v", err)
	}

	err = sendMessage(data, conn)
	if err != nil {
		return fmt.Errorf("WebSocket write message failed: %v", err)
	}
	return nil
}

func sendMessage(message []byte, dist *websocket.Conn) error {
	mu.Lock()
	defer mu.Unlock()

	err := dist.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	return nil
}

func readMessage(senderID int, conn *websocket.Conn) (int, []byte, error) {
	mu.Lock()
	defer mu.Unlock()
	log.Println(senderID)

	dataType, data, err := conn.ReadMessage()
	if err != nil {
		return 0, nil, fmt.Errorf("error read message: %v", err)
	}
	return dataType, data, nil
}
