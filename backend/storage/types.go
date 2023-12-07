package storage

type User struct {
	ID           string
	Email        string
	Password     string
	LastOnlineTS int64
	Bio          string
	RegisteredTS int64
}
