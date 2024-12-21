package models

type User struct {
	ID        int64
	Login     string
	Pass_hash []byte
}
