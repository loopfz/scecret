package models

import (
	"errors"
	"fmt"
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

// LocationCard represents the cards contained in a location.
// Decoupled from Location to allow things like multiple "A" locations
type LocationCard struct {
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
func LoadLocationFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*Location, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list locations")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"location"`).Where(
		squirrel.Eq{`id`: ID},
	)

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

// Create a link between a card and a location.
func (loc *Location) CreateLocationCard(db *gorp.DbMap, scenar *Scenario, letter string) (*LocationCard, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to create location card")
	}

	letter = strings.ToUpper(strings.TrimSpace(letter))

	cardDesc := fmt.Sprintf("%s - %s", loc.Name, letter)
	card, err := CreateCard(db, scenar, 0, cardDesc, &CardFace{}, &CardFace{})
	if err != nil {
		return nil, err
	}

	lc := &LocationCard{
		IDCard:     card.ID,
		IDLocation: loc.ID,
		Letter:     letter,
	}

	err = lc.Valid()
	if err != nil {
		return nil, err // TODO Tx
	}

	err = db.Insert(lc)
	if err != nil {
		return nil, err // TODO Tx
	}

	return lc, nil
}

// List a location's cards.
func (loc *Location) ListLocationCards(db *gorp.DbMap) ([]*LocationCard, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list location cards")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"location_card"`).Where(
		squirrel.Eq{`id_location`: loc.ID},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var lc []*LocationCard

	_, err = db.Select(&lc, query, args...)
	if err != nil {
		return nil, err
	}

	return lc, nil
}

// Load a location card,
func (loc *Location) LoadLocationCardFromID(db *gorp.DbMap, ID int64) (*LocationCard, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load location card")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"location_card"`).Where(
		squirrel.And{
			squirrel.Eq{`id`: ID},
			squirrel.Eq{`id_location`: loc.ID},
		},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var lc LocationCard

	_, err = db.Select(&lc, query, args...)
	if err != nil {
		return nil, err
	}

	return &lc, nil
}

// Update a location card.
func (lc *LocationCard) Update(db *gorp.DbMap, letter string) error {
	if db == nil {
		return errors.New("Missing db parameter to update location card")
	}

	lc.Letter = strings.ToUpper(strings.TrimSpace(letter))

	err := lc.Valid()
	if err != nil {
		return err
	}

	// Update card description with new letter
	loc, err := LoadLocationFromID(db, nil, lc.IDLocation)
	if err != nil {
		return err
	}
	card, err := LoadCardFromID(db, nil, lc.IDCard)
	if err != nil {
		return err
	}
	cardDesc := fmt.Sprintf("%s - %s", loc.Name, lc.Letter)
	err = card.Update(db, card.Number, cardDesc, card.Front, card.Back)
	if err != nil {
		return err
	}

	rows, err := db.Update(lc)
	if err != nil {
		return err // TODO Tx
	}
	if rows == 0 {
		return errors.New("No such location card to update") // TODO Tx
	}

	return nil
}

// Delete a location card.
func (lc *LocationCard) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to update location card")
	}

	// Delete card object
	card, err := LoadCardFromID(db, nil, lc.IDCard)
	if err != nil {
		return err
	}
	err = card.Delete(db)
	if err != nil {
		return err
	}

	rows, err := db.Delete(lc)
	if err != nil {
		return err // TODO Tx
	}
	if rows == 0 {
		return errors.New("No such location card to delete") // TODO Tx
	}

	return nil
}

// Verify that a LocationCard is valid before creating/updating it.
func (lc *LocationCard) Valid() error {
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	ok := false
	for _, l := range letters {
		if lc.Letter == l {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("Invalid location letter: %s", lc.Letter)
	}
	return nil
}

// TEMP ?
func (loc *Location) GetCards(db *gorp.DbMap) ([]*Card, error) {

	var c []*Card

	_, err := db.Select(&c, `SELECT card.* from "card" JOIN location_card ON card.id = location_card.id_card WHERE location_card.id_location = $1`, loc.ID)
	if err != nil {
		return nil, err
	}

	return c, nil
}
