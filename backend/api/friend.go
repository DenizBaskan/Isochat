package api

import (
	"backend/storage"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Request(c *fiber.Ctx) error {
	var body struct {
		Key string `json:"key"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	msg := "NEW_PUB_KEY:" + body.Key

	if err := storage.DB.Insert("USER_UPDATE (ID, MSG, TS) VALUES (?, ?, ?)", uuid.New().String(), msg, time.Now().Unix()); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}

func Accept(c *fiber.Ctx) error {
	var body struct {
		Key string `json:"key"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	msg := "NEW_PUB_KEY:" + body.Key

	if err := storage.DB.Insert("USER_UPDATE (ID, MSG, TS) VALUES (?, ?, ?)", uuid.New().String(), msg, time.Now().Unix()); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusCreated)
}
