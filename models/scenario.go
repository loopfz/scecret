package models

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
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

// List scenarios, optionally filtered by author.
func ListScenarios(db *gorp.DbMap, author *User) ([]*Scenario, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list scenarios")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"scenario"`)

	if author != nil {
		selector = selector.Where(
			squirrel.Eq{`id_author`: author.ID},
		)
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var s []*Scenario

	_, err = db.Select(&s, query, args...)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Load a scenario by id, optionally filtered by author.
func LoadScenarioFromID(db *gorp.DbMap, author *User, ID int64) (*Scenario, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load scenario")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"scenario"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if author != nil {
		selector = selector.Where(
			squirrel.Eq{`id_author`: author.ID},
		)
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var sc Scenario

	err = db.SelectOne(&sc, query, args...)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

// Update a scenario.
func (sc *Scenario) Update(db *gorp.DbMap, name string) error {
	if db == nil {
		return errors.New("Missing db parameter to update scenario")
	}

	sc.Name = name

	err := sc.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(sc)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such scenario to update")
	}

	return nil
}

// Delete a scenario.
func (sc *Scenario) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete scenario")
	}

	rows, err := db.Delete(sc)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such scenario to update")
	}

	return nil
}

// Verify that a scenario is valid before creating/updating it.
func (sc *Scenario) Valid() error {
	if sc.Name == "" {
		return errors.New("Empty name for scenario")
	}
	return nil
}
