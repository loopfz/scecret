package models

type User struct {
	ID           int64  `json:"-" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"`
	PasswordSalt string `json:"-" db:"password_salt"`
}

// TODO model functions
