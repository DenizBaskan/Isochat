package api

import (
	"os"
	"time"
	"backend/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	head := c.Request().Header.Cookie("auth_token")
	if head == nil {
		// authorization not provided
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var id string

	if err := storage.DB.GetRow([]any{&id}, "USER_ID FROM SESSION WHERE ID = ?", string(head)); err != nil {
		if err == sql.ErrNoRows {
			// auth header does not exist in db
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return err
	}

	// authorization header is valid
	c.Locals("ID", id)
	c.Locals("Session", string(head))

	return c.Next()
}

func newAuthCookie(val string) *fiber.Cookie {
	return &fiber.Cookie{
		Name: "auth_token",
		Value: val,
		Path: "/",
		Domain: os.Getenv("DOMAIN"),
		MaxAge: 60 * 24 * 365,
		Expires: time.Unix(time.Now().Unix() + (60 * 24 * 365), 0),
		Secure: true,
		HTTPOnly: false,
		SameSite: "Lax",
		SessionOnly: false,
	}
}
