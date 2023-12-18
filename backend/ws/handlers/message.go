package handlers

import (
	"backend/storage"
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

func HandleSendMessage(c *websocket.Conn, recv []byte, conns *map[string]*websocket.Conn, status int8) error {
	var req struct {
		Data struct {
			RecipientID string `json:"recipient_id"`
			Message     string `json:"message"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	if req.Data.Message == "" {
		return c.WriteMessage(websocket.TextMessage, reason("Message is empty"))
	} else if len(req.Data.Message) >= 500 {
		return c.WriteMessage(websocket.TextMessage, reason("Message cannot be longer than 500 characters"))
	}

	exists, err := storage.DB.Exists("USER WHERE ID = ?", req.Data.RecipientID)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(websocket.TextMessage, reason("User does not exist"))
	}
	
	if err := storage.DB.Insert("MESSAGE (ID, SENDER_ID, RECIPIENT_ID, DATA, SENT_TS, IS_READ) VALUES (?, ?, ?, ?, ?, ?)", 
	uuid.NewString(),
	c.Locals("ID"),
	req.Data.RecipientID,
	req.Data.Message,
	time.Now().Unix(),
	false); err != nil {
		return err
	}

	var resp struct {
		Updates struct { 
			Message struct {
				Data string	          `json:"data"`
				IsSender bool         `json:"is_sender"`
				SenderUsername string `json:"sender_username"`
			} `json:"message"`
		} `json:"updates"`
	}

	var username string

	if err := storage.DB.GetRow([]any{&username}, "USERNAME FROM USER WHERE ID = ?", c.Locals("ID")); err != nil {
		return err
	}

	// broadcast message to recipient
	if conn, ok := (*conns)[req.Data.RecipientID]; ok {
		cpy := resp

		cpy.Updates.Message.Data = req.Data.Message
		cpy.Updates.Message.IsSender = false
		cpy.Updates.Message.SenderUsername = username

		conn.WriteMessage(websocket.TextMessage, success(cpy, status))
	}

	resp.Updates.Message.Data = req.Data.Message
	resp.Updates.Message.IsSender = true
	resp.Updates.Message.SenderUsername = username

	return c.WriteMessage(websocket.TextMessage, success(resp, status))
}
