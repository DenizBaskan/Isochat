package api

import (
	"time"
	"backend/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

// must be used with auth middleware
func Authed(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}

func Auth(c *fiber.Ctx) error {
	head := c.Request().Header.Cookie("auth_token")
	if head == nil {
		// authorization not provided
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var (
		id     string
		expiry int
	)

	if err := storage.DB.GetRow([]any{&id, &expiry}, "USER_ID, EXPIRY FROM SESSION WHERE ID = ?", string(head)); err != nil {
		if err == sql.ErrNoRows {
			// auth header does not exist in db
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return err
	}

	if int(time.Now().Unix()) > expiry {
		// session expired
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// authorization header is valid
	c.Locals("ID", id)
	c.Locals("Session", string(head))

	return c.Next()
}

func newAuthCookie(val string, sec int64) *fiber.Cookie {
	return &fiber.Cookie{
		Name: "auth_token",
		Value: val,
		Path: "/",
		Expires: time.Unix(sec, 0),
		HTTPOnly: true,
	}
}
