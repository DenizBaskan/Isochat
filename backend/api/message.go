package api

import (
	"backend/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func HandleGetMessages(c *fiber.Ctx) error {
	rows, err := storage.DB.GetRows("DATA, SENDER_ID, USER.USERNAME FROM MESSAGE INNER JOIN USER ON SENDER_ID = USER.ID WHERE RECIPIENT_ID = ? AND SENDER_ID = ? OR RECIPIENT_ID = ? AND SENDER_ID = ? ORDER BY SENT_TS", c.Locals("ID"), c.Params("id"), c.Params("id"), c.Locals("ID"))
	if err != nil {
		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusOK)
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

	return c.Status(fiber.StatusOK).JSON(resp)
}
