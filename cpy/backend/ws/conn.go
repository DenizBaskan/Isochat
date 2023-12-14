package ws

import (
	"time"
	"backend/storage"
	"backend/ws/handlers"
	"encoding/json"
	"errors"

	"github.com/gofiber/contrib/websocket"
)

const (
	sendMessage = iota
	getMessages
	getFriends
	sendFriendRequest 
	acceptFriendRequest
	declineFriendRequest
	removeFriend
	removeFriendRequest
)

func Conn(c *websocket.Conn) {
	for {
		var msg struct {
			Status int8 `json:"status"`
		}

		mt, recv, err := c.ReadMessage()
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
			err = handlers.HandleSendMessage(c, mt, recv, msg.Status)
		case getMessages:
			err = handlers.HandleGetMessages(c, mt, recv, msg.Status)
		case getFriends:
			err = handlers.HandleGetFriends(c, mt, recv, msg.Status)
		case sendFriendRequest:
			err = handlers.HandleSendFriendRequest(c, mt, recv, msg.Status)
		case acceptFriendRequest:
			err = handlers.HandleAcceptFriendRequest(c, mt, recv, msg.Status)
		case declineFriendRequest:
			err = handlers.HandleDeclineFriendRequest(c, mt, recv, msg.Status)
		case removeFriend:
			err = handlers.HandleRemoveFriend(c, mt, recv, msg.Status)
		case removeFriendRequest:
			err = handlers.HandleRemoveFriendRequest(c, mt, recv, msg.Status)

		default:
			err = errors.New("invalid code recieved")
		}

		if err != nil {
			panic(err)
		}
	}
}
