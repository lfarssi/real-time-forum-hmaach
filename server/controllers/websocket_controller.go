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

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		utils.JSONResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer conn.Close()

	userID := r.Context().Value("user_id").(int)

	mu.Lock()
	ConnectedUsers[userID] = &Connection{Conn: conn, UserID: userID}
	mu.Unlock()

	broadcastUserList()

	for {
		err := handleChat(userID, conn)
		if err != nil {
			log.Println(err)
			break
		}
	}

	mu.Lock()
	delete(ConnectedUsers, userID)
	mu.Unlock()

	broadcastUserList()
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

func broadcastUserList() {
	mu.Lock()
	defer mu.Unlock()

	var userIDs []int
	for userID := range ConnectedUsers {
		userIDs = append(userIDs, userID)
	}
	message, err := json.Marshal(userIDs)
	if err != nil {
		log.Println("Error marshalling user IDs:", err)
		return
	}

	for _, connection := range ConnectedUsers {
		err := connection.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error broadcasting user list:", err)
		}
	}
}
