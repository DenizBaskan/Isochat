package api

import (
	"backend/storage"
	"crypto/rand"
	"crypto/tls"
	"net/mail"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// utility function to check if username is valid
func unameValid(user string) bool {
	if len(user) == 0 {
		return false
	}

	whitelisted := "0123456789abcdefghijklmnopqrstuvwxyz_-."

	for _, c := range user {
		if !strings.Contains(whitelisted, strings.ToLower(string(c))) {
			return false
		}
	}

	return true
}

// utility function for generating a six digit code
func genCode() (string, error) {
    codes := make([]byte, 6)
    if _, err := rand.Read(codes); err != nil {
        return "", err
    }

    for i := 0; i < 6; i++ {
        codes[i] = uint8(48 + (codes[i] % 10))
    }

    return string(codes), nil
}

// utility function to send a verification email
func mustSendEmailCode(email, code string) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		panic(err)
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USER"), os.Getenv("SMTP_PASS"))
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_USER"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your confirmation code")
	m.SetBody("text/html", "Your confirmation code is " + code)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

// utility function for checking the validity of an email domain
func isDomainAllowed(addr string) bool {
	// list of popular email domains
	whitelisted := []string{"gmail.com", "yahoo.com", "hotmail.com", "aol.com", "hotmail.co.uk", "hotmail.fr", "msn.com", "yahoo.fr", "wanadoo.fr", "orange.fr", "comcast.net", "yahoo.co.uk", "yahoo.com.br", "yahoo.co.in", "live.com", "rediffmail.com", "free.fr", "gmx.de", "web.de", "yandex.ru", "ymail.com", "libero.it", "outlook.com", "uol.com.br", "bol.com.br", "mail.ru", "cox.net", "hotmail.it", "sbcglobal.net", "sfr.fr", "live.fr", "verizon.net", "live.co.uk", "googlemail.com", "yahoo.es", "ig.com.br", "live.nl", "bigpond.com", "terra.com.br", "yahoo.it", "neuf.fr", "yahoo.de", "alice.it", "rocketmail.com", "att.net", "laposte.net", "facebook.com", "bellsouth.net", "yahoo.in", "hotmail.es", "charter.net", "yahoo.ca", "yahoo.com.au", "rambler.ru", "hotmail.de", "tiscali.it", "shaw.ca", "yahoo.co.jp", "sky.com", "earthlink.net", "optonline.net", "freenet.de", "t-online.de", "aliceadsl.fr", "virgilio.it", "home.nl", "qq.com", "telenet.be", "me.com", "yahoo.com.ar", "tiscali.co.uk", "yahoo.com.mx", "voila.fr", "gmx.net", "mail.com", "planet.nl", "tin.it", "live.it", "ntlworld.com", "arcor.de", "yahoo.co.id", "frontiernet.net", "hetnet.nl", "live.com.au", "yahoo.com.sg", "zonnet.nl", "club-internet.fr", "juno.com", "optusnet.com.au", "blueyonder.co.uk", "bluewin.ch", "skynet.be", "sympatico.ca", "windstream.net", "mac.com", "centurytel.net", "chello.nl", "live.ca", "aim.com", "bigpond.net.au"}
	
	for _, domain := range whitelisted {
		if strings.Contains(addr, "@" + domain) {
			return true
		}
	}

	return false
}

func HandlePostSignup(c *fiber.Ctx) error {
	var body struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		Code       string `json:"code"`
		Password   string `json:"password"`
		CaptchaKey string `json:"captcha_key"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}
	
	// retrieve code from cached storage
	code, err := storage.Cache.Get("email_address_" + body.Email)
	if err != nil {
		return err
	}

	if !(code != "" && code == body.Code) {
		// invalid email code provided
		return c.Status(fiber.StatusBadRequest).JSON(reason("Incorrect email code"))
	}
	
	if !unameValid(body.Username) {
		// username is invalid
		return c.Status(fiber.StatusBadRequest).JSON(reason("Username contains invalid characters"))
	}

	exists, err := storage.DB.Exists("USER WHERE USERNAME = ?", body.Username)
	if err != nil {
		return err
	}

	if exists {
		// user provided username that is already registered to another user
		return c.Status(fiber.StatusBadRequest).JSON(reason("Username is already registered"))
	}

	if len(body.Password) < 8 {
		// password is invalid, must be at least 8 chars
		return c.Status(fiber.StatusBadRequest).JSON(reason("Password must be at least eight characters"))
	}
	
	valid, err := isHCapValid(body.CaptchaKey)
	if err != nil {
		return err
	}

	if !valid {
		// user provided an invalid hcaptcha key
		return c.Status(fiber.StatusBadRequest).JSON(reason("Please solve the captcha"))
	}

	ts := time.Now().Unix()

	userID := uuid.New().String()
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
	if err != nil {
		return err
	}

	// create a new user in the db
	if err := storage.DB.Insert("USER (ID, EMAIL, USERNAME, PASSWORD, LAST_ONLINE_TS, REGISTERED_TS) VALUES (?, ?, ?, ?, ?, ?)",
		userID,
		body.Email,
		body.Username,
		pwdHash,
		ts,
		ts,
	); err != nil {
		return err
	}

	sess := uuid.New().String()

	if err := storage.DB.Insert("SESSION (ID, USER_ID) VALUES (?, ?)", sess, userID); err != nil {
		return err
	}

	c.Cookie(newAuthCookie(sess))

	return c.SendStatus(fiber.StatusCreated)
}

func HandlePostEmail(c *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		// user provided email is not a valid email
		return c.Status(fiber.StatusBadRequest).JSON(reason("Email is invalid"))
	}

	if !isDomainAllowed(body.Email) {
		// user provided a burner email
		return c.Status(fiber.StatusBadRequest).JSON(reason("We only allow registration from popular email domains"))
	}

	exists, err := storage.DB.Exists("USER WHERE EMAIL = ?", body.Email)
	if err != nil {
		return err
	}

	if exists {
		// user provided email that is already registered to another user
		return c.Status(fiber.StatusBadRequest).JSON(reason("Email is already registered"))
	}

	// generate six digit code
	code, err := genCode()
	if err != nil {
		return err
	}

	// Save email/code to cache with expiration of five minutes
	if err := storage.Cache.Set("email_address_" + body.Email, code, time.Duration((60 * 5) * float64(time.Second))); err != nil {
		return err
	}

	// send code to email, in another goroutine as it takes a while
	go mustSendEmailCode(body.Email, code)

	return c.SendStatus(fiber.StatusOK)
}
