package ws

import (
	"time"
	"backend/storage"
	"backend/ws/handlers"
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/websocket"
)

const sendMessage = iota

var conns = make(map[string]*websocket.Conn)

func Conn(c *websocket.Conn) {
	conns[c.Locals("ID").(string)] = c
	
	defer func() {
		delete(conns, c.Locals("ID").(string))
	}()

	for {
		var msg struct {
			Status int8 `json:"status"`
		}

		_, recv, err := c.ReadMessage()
		if err != nil {
			panic(err)
		}
		
		if err := json.Unmarshal(recv, &msg); err != nil {
			panic(err)
		}

		exists, err := storage.DB.Exists("SESSION WHERE ID = ?", c.Locals("Session"))
		if err != nil {
			panic(err)
		}

		if !exists {
			panic("session does not exist")
		}

		var expiry int

		if err := storage.DB.GetRow([]any{&expiry}, "EXPIRY FROM SESSION WHERE ID = ?", c.Locals("Session")); err != nil {
			panic(err)
		}

		if int(time.Now().Unix()) > expiry {
			// session expired
			panic("session has expired")
		}

		switch (msg.Status) {
		case sendMessage:
			err = handlers.HandleSendMessage(c, recv, &conns, msg.Status)

		default:
			err = errors.New("invalid code recieved")
		}

		if err != nil {
			panic(err)
		}
	}
}
