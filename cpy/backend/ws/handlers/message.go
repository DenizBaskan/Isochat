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
	
	rows, err := storage.DB.GetRows("DATA, SENDER_ID, USER.USERNAME FROM MESSAGE INNER JOIN USER ON SENDER_ID = USER.ID WHERE RECIPIENT_ID = ? AND SENDER_ID = ? OR RECIPIENT_ID = ? AND SENDER_ID = ? ORDER BY SENT_TS", c.Locals("ID"), req.Data.RecipientID, req.Data.RecipientID, c.Locals("ID"))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.WriteMessage(mt, success(nil, status)) 
		}

		return err
	}

	type Message struct {
		Data           string `json:"data"`
		IsSender       bool   `json:"is_sender"`
		SenderUsername string `json:"sender_username"`
	}

	var resp []Message

	for rows.Next() {
		var (
			m         Message
			sender_id string
		)
		rows.Scan(&m.Data, &sender_id, &m.SenderUsername)
		m.IsSender = sender_id == c.Locals("ID")

		resp = append(resp, m)
	}

	return c.WriteMessage(mt, success(resp, status))
}
