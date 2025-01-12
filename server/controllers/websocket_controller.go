package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"forum/server/utils"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
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
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message failed:", err)
			break
		}

		err = conn.WriteMessage(messageType, data)
		if err != nil {
			log.Println("WebSocket write message failed:", err)
			break
		}
	}

	mu.Lock()
	delete(ConnectedUsers, userID)
	mu.Unlock()

	broadcastUserList()
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
