package models

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

// LocationLink represents a link between a card and a Location.
// The card REVEALS the location.
type LocationLink struct {
	ID         int64 `json:"id" db:"id"`
	IDScenario int64 `json:"-" db:"id_scenario"`
	IDCard     int64 `json:"id_card" db:"id_card"`
	IDLocation int64 `json:"id_location" db:"id_location"`
}

// Create a link between a card and a location.
func CreateLocationLink(db *gorp.DbMap, card *Card, loc *Location) (*LocationLink, error) {
	if db == nil || card == nil || loc == nil {
		return nil, errors.New("Missing parameters to create location link")
	}

	if !loc.Hidden {
		// If the location needs to be revealed by another card, it is hidden.
		err := loc.Update(db, loc.Name, true, loc.Notes)
		if err != nil {
			return nil, err
		}
	}

	ll := &LocationLink{
		IDScenario: card.IDScenario,
		IDCard:     card.ID,
		IDLocation: loc.ID,
	}

	err := db.Insert(ll)
	if err != nil {
		return nil, err // TODO Tx
	}

	return ll, nil
}

// List location links, with filters.
func ListLocationLinks(db *gorp.DbMap, scenar *Scenario, card *Card, loc *Location) ([]*LocationLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load location links")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"location_link"`)

	if scenar != nil {
		selector = selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}
	if card != nil {
		selector = selector.Where(squirrel.Eq{`id_card`: card.ID})
	}
	if loc != nil {
		selector = selector.Where(squirrel.Eq{`id_location`: loc.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var ll []*LocationLink

	_, err = db.Select(&ll, query, args...)
	if err != nil {
		return nil, err
	}

	return ll, nil
}

// Loads a location link by ID. Optionally filtered by scenario.
func LoadLocationLinkFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*LocationLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card links")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"location_link"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar != nil {
		selector = selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var ll LocationLink

	err = db.SelectOne(&ll, query, args...)
	if err != nil {
		return nil, err
	}

	return &ll, nil
}

// Delete a location link.
func (ll *LocationLink) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete location link")
	}

	rows, err := db.Delete(ll)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such location link to delete")
	}

	return nil
}
