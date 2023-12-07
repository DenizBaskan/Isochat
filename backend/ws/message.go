package ws

import (
	"encoding/json"
	"time"

	"backend/storage"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

/*
CREATE TABLE USER_UPDATE (
    ID CHAR(36) NOT NULL PRIMARY KEY UNIQUE,
    USER_ID CHAR(36) NOT NULL,
    FOREIGN KEY (USER_ID) REFERENCES USER(ID),
    MSG TEXT NOT NULL,
    TS INT NOT NULL
);
*/

type MessageData struct {
	RecipientID string `json:"recipient_id"`
	Message     string `json:"message"`
}

func Message(c *websocket.Conn, mt int, recv any) error {
	var data MessageData

	if err := json.Unmarshal(recv.([]byte), &data); err != nil {
		return err
	}

	msg := "MSG:" + data.Message

	if err := storage.DB.Insert("USER_UPDATE (ID, USER_ID, MSG, TS) VALUES (?, ?, ?, ?)", uuid.New().String(), data.RecipientID, msg, time.Now().Unix()); err != nil {
		return err
	}

	return c.WriteMessage(mt, []byte(`{"success":true}`))
}
