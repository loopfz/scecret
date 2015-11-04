package models

import (
	"errors"
	"strings"

	"github.com/go-gorp/gorp"
)

// Scenario is the highest-level object.
// All other game objects belong to a scenario.
// A scenario belongs to a user (author).
type Scenario struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	IDAuthor int64  `json:"-" db:"id_author"`
}

// Create a scenario.
func CreateScenario(db *gorp.DbMap, name string, author *User) (*Scenario, error) {
	if db == nil || author == nil {
		return nil, errors.New("Missing parameters to create scenario")
	}

	sc := &Scenario{
		Name:     strings.TrimSpace(name),
		IDAuthor: author.ID,
	}

	err := sc.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(sc)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

// Verify that a scenario is valid before creating/updating it.
func (sc *Scenario) Valid() error {
	if sc.Name == "" {
		return errors.New("Empty name for scenario")
	}
	return nil
}

// TODO model functions
