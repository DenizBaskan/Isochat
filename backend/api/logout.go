package api

import (
	"backend/storage"

	"github.com/gofiber/fiber/v2"
)

func HandlePostLogout(c *fiber.Ctx) error {
	exists, err := storage.DB.Exists("SESSION WHERE ID = ?", c.Locals("Session"))
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(reason("Session does not exist"))
	}

	if err := storage.DB.Delete("SESSION WHERE ID = ?", c.Locals("Session")); err != nil {
		return err
	}

	c.Cookie(newAuthCookie("", 0))

	return c.SendStatus(fiber.StatusOK)
}
