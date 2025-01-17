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

	broadcastUserList()
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

		err = models.SendMessage(message)
		if err != nil {
			log.Println("Failed to save message in database: ", err)
			utils.SendErrorMessage(conn, "Internal Server Srror")
			continue
		}

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

	receiver, exists := ConnectedUsers[message.ReceiverID]
	if !exists {
		return fmt.Errorf("not found")
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	// Send the message to the receiver
	err = receiver.Conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		// Close the connection and remove the user from the connected users map
		receiver.Conn.Close()
		delete(ConnectedUsers, message.ReceiverID)

		return fmt.Errorf("receiver disconnected: %v", err)
	}

	return nil
}

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

		err = connection.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// if there was an error it means that the user is disconnected
			connection.Conn.Close()
			delete(ConnectedUsers, userID)
		}
	}
}
