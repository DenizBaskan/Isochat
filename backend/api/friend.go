package api

import (
	"backend/storage"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func HandleSendFriendRequest(c *fiber.Ctx) error {
	var body struct {
		Username string `json:"username"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	exists, err := storage.DB.Exists("USER WHERE USERNAME = ?", body.Username)
	if err != nil {
		return err
	}

	// if user with the username of the friend request does not exist
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Username does not exist"))
	}

	// user id of the recipient
	var userID string

	if err = storage.DB.GetRow([]any{&userID}, "ID FROM USER WHERE USERNAME = ?", body.Username); err != nil {
		return err
	}

	// if user has tried to send a friend request to themselves
	if userID == c.Locals("ID") {
		return c.Status(fiber.StatusBadRequest).JSON(reason("You cannot send a friend request to yourself"))
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

			return c.SendStatus(fiber.StatusCreated)
		}

		return err
	}

	if accepted {
		return c.Status(fiber.StatusBadRequest).JSON(reason("You are already friend with this user"))
	}
	
	if senderID == userID {
		return c.Status(fiber.StatusBadRequest).JSON(reason("You have already sent this user a friend request"))
	}

	return c.Status(fiber.StatusBadRequest).JSON(reason("This user has already sent you a friend request"))
}

func HandleDeclineFriendRequest(c *fiber.Ctx) error {
	var body struct {
		FriendID string `json:"friend_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	// check if there is a friend request for the user matching the id provided
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? AND ID = ?", c.Locals("ID"), body.FriendID)
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Freind request does not exist"))
	}

	var accepted bool

	// get the accepted boolean for the friend record
	if err := storage.DB.GetRow([]any{&accepted}, "ACCEPTED FROM FRIEND WHERE ID = ?", body.FriendID); err != nil {
		return err
	}

	if accepted {
		return c.Status(fiber.StatusBadRequest).JSON(reason("You are already friend with this user"))
	}

	if err := storage.DB.Delete("FRIEND WHERE ID = ?", body.FriendID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func HandleAcceptFriendRequest(c *fiber.Ctx) error {
	var body struct {
		FriendID string `json:"friend_id"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	// check if there is a friend request for the user matching the id provided
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? AND ID = ?", c.Locals("ID"), body.FriendID)
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Friend request does not exist"))
	}

	var accepted bool

	// get the accepted boolean for the friend record
	if err := storage.DB.GetRow([]any{&accepted}, "ACCEPTED FROM FRIEND WHERE ID = ?", body.FriendID); err != nil {
		return err
	}

	if accepted {
		return c.Status(fiber.StatusBadRequest).JSON(reason("You are already friend with this user"))
	}

	if err := storage.DB.Update("FRIEND SET ACCEPTED = ?, ACCEPTED_TS = ? WHERE ID = ?", true, time.Now().Unix(), body.FriendID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func HandleGetFriends(c *fiber.Ctx) error {
	type Friend struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		UserID   string `json:"user_id"`
	}

	var resp struct {
		Pending  []Friend `json:"pending"`
		Incoming []Friend `json:"incoming"`
		Friends  []Friend `json:"friends"`
	}

	// query the friend rows involving the user
	rows, err := storage.DB.GetRows("FRIEND.ID, RECIPIENT_ID, SENDER_ID, ACCEPTED, U1.USERNAME, U2.USERNAME FROM FRIEND INNER JOIN USER U1 ON RECIPIENT_ID = U1.ID INNER JOIN USER U2 ON SENDER_ID = U2.ID WHERE SENDER_ID = ? OR RECIPIENT_ID = ?", c.Locals("ID"), c.Locals("ID"))
	if err != nil {
		if err == sql.ErrNoRows {
			// user has no friend records yet
			return c.SendStatus(fiber.StatusOK)
		}

		return err
	}

	for rows.Next() {
		var  (
			recipientID, senderID, recipientUsername, senderUsername string
			accepted                                                     bool
		)

		var f Friend

		// set all the variables from the sql query
		rows.Scan(&f.ID, &recipientID, &senderID, &accepted, &recipientUsername, &senderUsername)

		f.UserID = senderID
		f.Username = senderUsername

		// sender is user
		if senderID == c.Locals("ID") {
			f.UserID = recipientID
			f.Username = recipientUsername
		}
		
		// accepted
		if accepted {
			// if the friend record is accepted, append the newly created friend object to the response data
			resp.Friends = append(resp.Friends, f)
		} else {
			if senderID == c.Locals("ID") { // if sender id is the user
				resp.Pending = append(resp.Pending, f)
			} else {
				resp.Incoming = append(resp.Incoming, f)
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func HandleRemoveFriend(c *fiber.Ctx) error {
	// check if there is a friend record involving the user where the "accepted" field is true
	exists, err := storage.DB.Exists("FRIEND WHERE RECIPIENT_ID = ? OR SENDER_ID = ? AND ID = ? AND ACCEPTED = ?", c.Locals("ID"), c.Locals("ID"), c.Params("id"), true)
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Friend does not exist"))
	}

	// remove the friend record
	if err := storage.DB.Delete("FRIEND WHERE ID = ?", c.Params("id")); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func HandleRemoveFriendRequest(c *fiber.Ctx) error {
	// check if there is a friend record to the user which has not yet been accepted
	exists, err := storage.DB.Exists("FRIEND WHERE SENDER_ID = ? AND ACCEPTED = ? AND ID = ?", c.Locals("ID"), false, c.Params("id"))
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Friend does not exist"))
	}

	// remove the friend record
	if err := storage.DB.Delete("FRIEND WHERE ID = ?", c.Params("id")); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
