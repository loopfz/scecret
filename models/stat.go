package models

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

type Stat struct {
	ID          int64  `json:"id" db:"id"`
	IDScenario  int64  `json:"-" db:"id_scenario"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IDIcon      int64  `json:"id_icon" db:"id_icon"`
}

// Create a stat object.
func CreateStat(db *gorp.DbMap, scenar *Scenario, ico *Icon, name string, description string) (*Stat, error) {
	if db == nil || scenar == nil || ico == nil {
		return nil, errors.New("Missing parameters to create stat")
	}

	st := &Stat{
		IDScenario:  scenar.ID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		IDIcon:      ico.ID,
	}

	err := st.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(st)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// List stats, optionally filtered by scenario.
func ListStats(db *gorp.DbMap, scenar *Scenario) ([]*Stat, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list stats")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"stat"`)

	if scenar != nil {
		selector = selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var st []*Stat

	_, err = db.Select(&st, query, args...)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// Load stat by id, optionally filtered by scenario.
func LoadStatFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*Stat, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load stat")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"stat"`)

	if scenar != nil {
		selector = selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var st Stat

	err = db.SelectOne(&st, query, args...)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

// Update stat object.
func (st *Stat) Update(db *gorp.DbMap, ico *Icon, name string, description string) error {
	if db == nil || ico == nil {
		return errors.New("Missing parameters to update stat")
	}

	st.IDIcon = ico.ID
	st.Name = name
	st.Description = description

	err := st.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(st)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such stat to update")
	}

	return nil
}

// Delete a stat object.
func (st *Stat) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete stat")
	}

	rows, err := db.Delete(st)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such stat to delete")
	}

	return nil
}

// Verify that a stat object is valid before creating/updating it.
func (st *Stat) Valid() error {
	if st.Name == "" {
		return errors.New("Empty name")
	}
	return nil
}
