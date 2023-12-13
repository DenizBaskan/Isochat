package handlers

import (
	"backend/storage"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

func HandleSendFriendRequest(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			Username string `json:"username"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	exists, err := storage.DB.Exists("USER WHERE USERNAME = ?", req.Data.Username)
	if err != nil {
		return err
	}

	// if user with the username of the friend request does not exist
	if !exists {
		return c.WriteMessage(mt, reason("Username does not exist"))
	}

	// user id of the recipient
	var userID string

	if err = storage.DB.GetRow([]any{&userID}, "ID FROM USER WHERE USERNAME = ?", req.Data.Username); err != nil {
		return err
	}

	// if user has tried to send a friend request to themselves
	if userID == c.Locals("ID") {
		return c.WriteMessage(mt, reason("You cannot send a friend request to yourself"))
	}

	var (
		senderID string
		accepted bool
	)

	// get a row of the friend relationship between the two
	if err = storage.DB.GetRow([]any{&accepted, &senderID}, "ACCEPTED, SENDER_ID FROM FRIEND WHERE SENDER_ID = ? AND RECIPIENT_ID = ? OR SENDER_ID = ? AND RECIPIENT_ID = ?", c.Locals("ID"), userID, userID, c.Locals("ID")); err != nil {
		// if there is not a friend record between the two, create one
		if err == sql.ErrNoRows {
			if err := storage.DB.Insert("FRIEND (ID, SENDER_ID, RECIPIENT_ID, ACCEPTED, SENT_TS) VALUES (?, ?, ?, ?, ?)", 
			uuid.NewString(),
			c.Locals("ID"),
			userID,
			false,
			time.Now().Unix()); err != nil {
				return err
			}

			return c.WriteMessage(mt, success(nil, status))
		}

		return err
	}

	if accepted {
		return c.WriteMessage(mt, reason("You are already friends with this user"))
	}
	
	if senderID == userID {
		return c.WriteMessage(mt, reason("You have already sent this user a friend request"))
	}

	return c.WriteMessage(mt, reason("This user has already sent you a friend request"))
}

func HandleDeclineFriendRequest(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			FriendID string `json:"friend_id"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	// check if there is a friend request for the user matching the id provided
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? AND ID = ?", c.Locals("ID"), req.Data.FriendID)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(mt, reason("Friend request does not exist"))
	}

	var accepted bool

	// get the accepted boolean for the friend record
	if err := storage.DB.GetRow([]any{&accepted}, "ACCEPTED FROM FRIEND WHERE ID = ?", req.Data.FriendID); err != nil {
		return err
	}

	if accepted {
		return c.WriteMessage(mt, reason("You are already friends with this user"))
	}

	if err := storage.DB.Delete("FRIEND WHERE ID = ?", req.Data.FriendID); err != nil {
		return err
	}

	return c.WriteMessage(mt, success(nil, status))
}

func HandleAcceptFriendRequest(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			FriendID string `json:"friend_id"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	// check if there is a friend request for the user matching the id provided
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? AND ID = ?", c.Locals("ID"), req.Data.FriendID)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(mt, reason("Friend request does not exist"))
	}

	var accepted bool

	// get the accepted boolean for the friend record
	if err := storage.DB.GetRow([]any{&accepted}, "ACCEPTED FROM FRIEND WHERE ID = ?", req.Data.FriendID); err != nil {
		return err
	}

	if accepted {
		return c.WriteMessage(mt, reason("You are already friends with this user"))
	}

	if err := storage.DB.Update("FRIEND SET ACCEPTED = ?, ACCEPTED_TS = ? WHERE ID = ?", true, time.Now().Unix(), req.Data.FriendID); err != nil {
		return err
	}

	return c.WriteMessage(mt, success(nil, status))
}

func HandleGetFriends(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var (
		id, recipientID, senderID, recpientName, senderName string
		accepted                                            bool
	)

	type (
		Friend struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			UserID   string `json:"user_id"`
		}
	
		Request struct {
			ID                string `json:"id"`
			SenderID          string `json:"sender_id"`
			RecipientID       string `json:"recpient_id"`
			SenderUsername    string `json:"sender_username"`
			RecipientUsername string `json:"recipient_username"`
			Accepted          bool   `json:"accepted"`
		}
	)

	var resp struct {
		Pending  []Request `json:"pending"`
		Incoming []Request `json:"incoming"`
		Friends  []Friend  `json:"friends"`
	}

	// query the friend rows involving the user
	rows, err := storage.DB.GetRows([]any{&id, &recipientID, &senderID, &accepted, &recpientName, &senderName}, "FRIEND.ID, RECIPIENT_ID, SENDER_ID, ACCEPTED, U1.USERNAME, U2.USERNAME FROM FRIEND INNER JOIN USER U1 ON RECIPIENT_ID = U1.ID INNER JOIN USER U2 ON SENDER_ID = U2.ID WHERE SENDER_ID = ? OR RECIPIENT_ID = ?", c.Locals("ID"), c.Locals("ID"))
	if err != nil {
		if err == sql.ErrNoRows {
			// user has no friend records yet
			return c.WriteMessage(mt, success(nil, status)) 
		}

		return err
	}

	for _, row := range rows {
		var r Request
		
		// set all the variables from the sql query
		r.ID = *row[0].(*string)
		r.RecipientID = *row[1].(*string)
		r.SenderID = *row[2].(*string)
		r.Accepted = *row[3].(*bool)
		r.RecipientUsername = *row[4].(*string)
		r.SenderUsername = *row[5].(*string)

		// accepted
		if r.Accepted {
			// if the friend record is accepted, append the newly created friend object to the response data
			var f Friend
			f.ID = r.ID
			f.UserID = r.SenderID
			f.Username = r.SenderUsername

			// sender is user
			if r.SenderID == c.Locals("ID") {
				f.UserID = r.RecipientID
				f.Username = r.RecipientUsername
			}
			
			resp.Friends = append(resp.Friends, f)
		} else {
			if r.SenderID == c.Locals("ID") { // if sender id is the user
				resp.Pending = append(resp.Pending, r)
			} else {
				resp.Incoming = append(resp.Pending, r)
			}
		}
	}

	return c.WriteMessage(mt, success(resp, status))
}

func HandleRemoveFriend(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			FriendID string `json:"friend_id"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	// check if there is a friend record involving the user where the "accepted" field is true
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? OR SENDER_ID = ? AND ID = ? AND ACCEPTED = ?", c.Locals("ID"), c.Locals("ID"), req.Data.FriendID, true)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(mt, reason("Friend does not exist"))
	}

	// remove the friend record
	if err := storage.DB.Delete("FRIEND WHERE ID = ?", req.Data.FriendID); err != nil {
		return err
	}

	return c.WriteMessage(mt, success(nil, status))
}

func HandleRemoveFriendRequest(c *websocket.Conn, mt int, recv []byte, status int8) error {
	var req struct {
		Data struct {
			FriendID string `json:"friend_id"`
		}
	}

	if err := json.Unmarshal(recv, &req); err != nil {
		return err
	}

	// check if there is a friend record to the user which has not yet been accepted
	exists, err := storage.DB.Exists("FRIEND WHERE SENDER_ID = ? AND ACCEPTED = ? AND ID = ?", c.Locals("ID"), false, req.Data.FriendID)
	if err != nil {
		return err
	}

	if !exists {
		return c.WriteMessage(mt, reason("Friend does not exist"))
	}

	// remove the friend record
	if err := storage.DB.Delete("FRIEND WHERE ID = ?", req.Data.FriendID); err != nil {
		return err
	}

	return c.WriteMessage(mt, success(nil, status))
}
