package api

import (
	"backend/storage"

	"github.com/gofiber/fiber/v2"
)

func HandlePostLogout(c *fiber.Ctx) error {
	// delete user session
	if err := storage.DB.Delete("FROM SESSION WHERE ID = ?", c.Get("Authorization")); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusOK)
}
