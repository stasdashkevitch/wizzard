package entity

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}
