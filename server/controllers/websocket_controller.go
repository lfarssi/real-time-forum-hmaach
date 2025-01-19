package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"forum/server/models"
	"forum/server/utils"
	"forum/server/validators"

	"github.com/gorilla/websocket"
)

var (
	ConnectedUsers = make(map[int]*websocket.Conn)
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	mu sync.Mutex
)

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
	var connection *websocket.Conn = conn

	mu.Lock()
	ConnectedUsers[userID] = connection
	mu.Unlock()

	broadcastOnlineUserList()

	for {
		err = handleChat(userID, conn)
		if err != nil {
			// if the user disconnected
			break
		}
	}

	// Clean up on disconnection
	mu.Lock()
	delete(ConnectedUsers, userID)
	mu.Unlock()

	broadcastOnlineUserList()
}

func handleChat(userID int, conn *websocket.Conn) error {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			utils.SendErrorMessage(conn, "Internal Server Srror")
			return fmt.Errorf("WebSocket read message failed: %v", err)
		}

		// Validate the chat message
		message, err := validators.ChatMessageRequest(userID, data)
		if err != nil {
			utils.SendErrorMessage(conn, err.Error())
			continue
		}

		// Get sender information
		sender, err := models.GetUserInfo(userID)
		if err != nil {
			log.Printf("Failed to get sender info for user %d: %v\n", userID, err)
			utils.SendErrorMessage(conn, "Internal Server Srror")
			continue
		}

		message.SenderID = sender.ID
		message.Sender = sender.Nickname

		err = models.StoreMessage(message)
		if err != nil {
			log.Println("Failed to save message in database: ", err)
			utils.SendErrorMessage(conn, "Internal Server Srror")
			continue
		}

		broadcastMessage("refresh-users")

		// Send the message to the receiver
		err = sendMessage(message)
		if err != nil {
			if err.Error() == "not found" {
				continue
			}
			log.Printf("Error sending message to receiver: %v\n", err)
			utils.SendErrorMessage(conn, "Failed to send message")
			continue
		}
	}
}

func sendMessage(message models.Message) error {
	mu.Lock()
	defer mu.Unlock()

	receiverConn, exists := ConnectedUsers[message.ReceiverID]
	if !exists {
		return fmt.Errorf("not found")
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	// Send the message to the receiver
	err = receiverConn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		// Close the connection and remove the user from the connected users map
		receiverConn.Close()
		delete(ConnectedUsers, message.ReceiverID)

		return fmt.Errorf("receiver disconnected: %v", err)
	}

	return nil
}

func broadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	// Broadcast to all connections
	for userID, connection := range ConnectedUsers {
		message, err := json.Marshal(map[string]interface{}{
			"type": message,
		})
		if err != nil {
			log.Printf("Error marshalling user list for user %d: %v\n", userID, err)
			continue
		}

		err = connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			connection.Close()
			delete(ConnectedUsers, userID)
		}
	}
}

func broadcastOnlineUserList() {
	mu.Lock()
	defer mu.Unlock()

	// Prepare a list of all user IDs
	userIDs := make([]int, 0, len(ConnectedUsers))
	for userID := range ConnectedUsers {
		userIDs = append(userIDs, userID)
	}

	// Broadcast to all connections
	for userID, connection := range ConnectedUsers {
		// Use a single loop to prepare a filtered list excluding the current user's ID
		filteredUserIDs := make([]int, 0, len(userIDs)-1)
		for _, id := range userIDs {
			if id != userID {
				filteredUserIDs = append(filteredUserIDs, id)
			}
		}

		// Marshal the filtered list into JSON format and send it to the current connection
		data := map[string]interface{}{
			"type":  "users-status",
			"users": filteredUserIDs,
		}
		message, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling user list for user %d: %v\n", userID, err)
			continue
		}

		err = connection.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// if there was an error it means that the user is disconnected
			connection.Close()
			delete(ConnectedUsers, userID)
		}
	}
}
