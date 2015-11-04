package models

import (
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

// Location represents a location the players can visit.
// It is composed of several cards.
type Location struct {
	ID         int64  `json:"id" db:"id"`
	IDScenario int64  `json:"-" db:"id_scenario"`
	Name       string `json:"name" db:"name"`
	Hidden     bool   `json:"json:"hidden" db:"hidden"`
}

// LocationCards represents the cards contained in a location.
// Decoupled from Location to allow things like multiple "A" locations
type LocationCards struct {
	ID         int64  `json:"id" db:"id"`
	IDLocation int64  `json:"id_location" db:"id_location"`
	IDCard     int64  `json:"id_card" db:"id_card"`
	Letter     string `json:"letter" db:"letter"`
}

// Create a location.
func CreateLocation(db *gorp.DbMap, scenar *Scenario, Name string, Hidden bool) (*Location, error) {
	if db == nil || scenar == nil {
		return nil, errors.New("Missing parameters to create location")
	}

	loc := &Location{
		IDScenario: scenar.ID,
		Name:       strings.TrimSpace(Name),
		Hidden:     Hidden,
	}

	err := loc.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(loc)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// List locations, with filters.
func ListLocations(db *gorp.DbMap, scenar *Scenario) ([]*Location, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list locations")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"location"`)

	if scenar != nil {
		selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var loc []*Location

	_, err = db.Select(&loc, query, args...)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

// Load a Location from ID. Optionally filtered by scenario.
func LoadLocationFromID(db *gorp.DbMap, scenar *Scenario) (*Location, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list locations")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"location"`)

	if scenar != nil {
		selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var loc Location

	err = db.SelectOne(&loc, query, args...)
	if err != nil {
		return nil, err
	}

	return &loc, nil
}

// Update a location.
func (loc *Location) Update(db *gorp.DbMap, Name string, Hidden bool) error {
	if db == nil {
		return errors.New("Missing db parameter to update location")
	}

	loc.Name = Name
	loc.Hidden = Hidden

	// TODO harmonize LocationCards descriptions (new name)

	err := loc.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(loc)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such location to update")
	}
	return nil
}

// Delete a location.
func (loc *Location) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete location")
	}

	rows, err := db.Delete(loc)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such location to delete")
	}
	return nil
}

// Verify that a Location is valid before creating/updating it.
func (loc *Location) Valid() error {
	if loc.Name == "" {
		return errors.New("Empty location name")
	}
	return nil
}
