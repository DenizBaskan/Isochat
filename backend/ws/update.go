package ws

import (
	//"encoding/json"

	"github.com/gofiber/contrib/websocket"
)

func UserUpdate(c *websocket.Conn, mt int, recv []byte) error {
	println(string(recv))

	return nil
}
