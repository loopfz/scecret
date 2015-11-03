package user

type User struct {
	ID           int64  `json:"-" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash `json:"-" db:"password_hash"`
	PasswordSalt `json:"-" db:"password_salt"`
}

// TODO model functions
