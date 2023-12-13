package handlers

import (
	"backend/storage"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

func HandleSendMessage(c *websocket.Conn, mt int, recv []byte, status int8) error {
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
		return c.WriteMessage(mt, reason("Message is empty"))
	}

	exists, err := storage.DB.Exists("USER WHERE ID = ?", req.Data.RecipientID)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(mt, reason("User does not exist"))
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

	return c.WriteMessage(mt, success(nil, status))
}

func HandleGetMessages(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			RecipientID string `json:"recipient_id"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	var (
		id, senderID, recpientID, data string
		sentTS                         int
		isRead                         bool
	)

	rows, err := storage.DB.GetRows([]any{&id, &senderID, &recpientID, &data, &sentTS, &isRead}, "* FROM MESSAGE WHERE RECIPIENT_ID = ? OR SENDER_ID = ? AND RECIPIENT_ID = ? OR SENDER_ID = ? ORDER BY SENT_TS", c.Locals("ID"), c.Locals("ID"), req.Data.RecipientID, req.Data.RecipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.WriteMessage(mt, success(nil, status)) 
		}
	}

	return c.WriteMessage(mt, success(rows, status))
}
