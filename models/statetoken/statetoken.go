package statetoken

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
)

// StateToken represents base game tokens with different icons
// that maintain game state (e.g. unlock access to cards)
// These are bootstrapped in DB, only need Load functions.
type StateToken struct {
	ID     int64 `json:"id" db:"id"`
	IDIcon int64 `json:"id_icon" db:"id_icon"`
}

// Loads a state token by ID
func LoadFromID(db *gorp.DbMap, ID int64) (*StateToken, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load state token")
	}

	query, args, err := squirrel.Select(`*`).From(`"state_token"`).Where(
		squirrel.Eq{`id`: ID},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var st StateToken

	err = db.SelectOne(&st, query, args...)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

// Loads all state tokens
func LoadAll(db *gorp.DbMap) ([]*StateToken, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load state tokens")
	}

	query, args, err := squirrel.Select(`*`).From(`"state_token"`).ToSql()

	if err != nil {
		return nil, err
	}

	var st []*StateToken

	_, err = db.Select(&st, query, args...)
	if err != nil {
		return nil, err
	}

	return st, nil
}
