package api

import (
	"backend/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

/*
CREATE TABLE USER (
	ID CHAR(36) NOT NULL PRIMARY KEY UNIQUE,
	EMAIL VARCHAR(320) NOT NULL UNIQUE,
	USERNAME CHAR(20) NOT NULL UNIQUE,
	PASSWORD CHAR(72) NOT NULL,
	LAST_ONLINE_TS INT NOT NULL,
    BIO CHAR(100),
	REGISTERED_TS INT NOT NULL
);
*/

func HandleGetFriends(c *fiber.Ctx) error {
	id := c.Query("username")

	var (

	)

	if err := storage.DB.GetRow([]any{&id}, "* FROM USER WHERE ID = ?", id); err != nil {
		if err == sql.ErrNoRows {
			// auth header does not exist in db
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return err
	}

	// authorization header is valid
	c.Locals("ID", id)

	return c.Next()
}
