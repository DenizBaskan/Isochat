package api

import (
	"backend/storage"
	"backend/ws"
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/websocket"
)

const (
	UserUpdate = iota
	SendMessage
)

type RecieveMessage struct {
	Code  int8   `json:"msg"`
	Token string `json:"token"`
	Data  any    `json:"Data"`
}

func WSConn(c *websocket.Conn) {
	if err := func() error {
		for {
			var msg RecieveMessage
	
			mt, recv, err := c.ReadMessage()
			if err != nil {
				return err
			}

			if c.Locals("Session") != msg.Token {
				return errors.New("session does not match old one")
			}
			
			if err := json.Unmarshal(recv, &msg); err != nil {
				return err
			}

			exists, err := storage.DB.Exists("SESSION WHERE ID = ?", msg.Token)
			if err != nil {
				return err
			}

			if !exists {
				return errors.New("session does not exist")
			}
	
			switch (msg.Code) {
			case UserUpdate:
				err = ws.UserUpdate(c, mt, recv)
			case SendMessage:
				err = ws.Message(c, mt, msg.Data)
	
			default:
				err = errors.New("invalid code recieved")
			}
	
			if err != nil {
				return err
			}
		}
	}; err != nil {
		panic(err)
	}
}
