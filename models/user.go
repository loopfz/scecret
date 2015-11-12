package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/scrypt"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/securerandom"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

const (
	SALT_LEN         = 32
	PASSWORD_MIN_LEN = 8
)

type User struct {
	ID           int64  `json:"-" db:"id"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"-" db:"password_hash"`
	PasswordSalt string `json:"-" db:"password_salt"`
	pwplain      string `json:"-" db:"-"`
}

// Create a user.
func CreateUser(db *gorp.DbMap, email string, password string) (*User, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to create user")
	}

	salt, err := securerandom.RandomString(SALT_LEN)
	if err != nil {
		return nil, err
	}

	u := &User{
		Email:        email,
		pwplain:      password,
		PasswordSalt: salt,
	}

	u.hash()

	err = u.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Load a user by email.
func LoadUserFromEmail(db *gorp.DbMap, email string) (*User, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load user")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"user"`).Where(
		squirrel.Eq{`email`: email},
	)

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var u User

	err = db.SelectOne(&u, query, args...)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// Update a user.
func (u *User) Update(db *gorp.DbMap, email string, password string) error {
	if db == nil {
		return errors.New("Missing db parameter to update user")
	}

	salt, err := securerandom.RandomString(SALT_LEN)
	if err != nil {
		return err
	}

	u.Email = email
	u.pwplain = password
	u.PasswordSalt = salt

	u.hash()

	err = u.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(u)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such user to update")
	}

	return nil
}

// Verify that a user is valid before creating/updating it.
func (u *User) Valid() error {
	if u.Email == "" {
		// TODO match email regex
		return errors.New("Empty email")
	}
	if len(u.pwplain) < PASSWORD_MIN_LEN {
		return fmt.Errorf("Password too short: min %d", PASSWORD_MIN_LEN)
	}
	if u.PasswordHash == "" || u.PasswordSalt == "" {
		return errors.New("Missing hashed password")
	}
	return nil
}

// Hash the user's password (scrypt).
func (u *User) hash() {

	key, err := scrypt.Key([]byte(u.pwplain), []byte(u.PasswordSalt), 16384, 8, 1, 32)
	if err == nil {
		u.PasswordHash = string(key)
	}
}

// Check if the user's password matches the plain parameter.
func (u *User) PasswordEquals(pw string) bool {
	u2 := &User{
		pwplain:      pw,
		PasswordSalt: u.PasswordSalt,
	}
	u2.hash()

	return u.PasswordHash == u2.PasswordHash
}
