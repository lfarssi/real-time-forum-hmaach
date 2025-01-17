package utils

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SendErrorMessage(conn *websocket.Conn, errorMessage string) {
	errorResponse := map[string]string{
		"type":    "error",
		"message": errorMessage,
	}

	messageJSON, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("Error marshalling error response: %v\n", err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, messageJSON)
}
