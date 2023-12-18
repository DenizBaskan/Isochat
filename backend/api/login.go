package api

import (
	"time"
	"backend/storage"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

func HandlePostLogin(c *fiber.Ctx) error {
	var body struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		CaptchaKey string `json:"captcha_key"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	valid, err := isHCapValid(body.CaptchaKey)
	if err != nil {
		return err
	}

	if !valid {
		// user provided an invalid hcaptcha key
		return c.Status(fiber.StatusBadRequest).JSON(reason("Please solve the captcha"))
	}

	var password, id string

	if err := storage.DB.GetRow([]any{&password, &id}, "PASSWORD, ID FROM USER WHERE EMAIL = ?", body.Email); err != nil {
		if err == sql.ErrNoRows {
			// email is not associated with a user
			return c.Status(fiber.StatusBadRequest).JSON(reason("Email does not exist"))
		}

		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(body.Password)); err != nil {
		// password hash did not match password meaning incorrect password
		return c.Status(fiber.StatusBadRequest).JSON(reason("Incorrect password"))
	}

	sess := uuid.New().String()

	// create a new session
	if err := storage.DB.Insert("SESSION (ID, USER_ID, EXPIRY) VALUES (?, ?, ?)", sess, id, time.Now().Unix() + (60 * 60 * 24 * 365)); err != nil {
		return err
	}

	c.Cookie(newAuthCookie(sess, time.Now().Unix() + (60 * 60 * 24 * 365)))

	return c.SendStatus(fiber.StatusOK)
}
