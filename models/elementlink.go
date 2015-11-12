package models

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

type ElementLink struct {
	ID         int64 `json:"id" db:"id"`
	IDScenario int64 `json:"id_scenario" db:"id_scenario"`
	IDElement  int64 `json:"id_element" db:"id_element"`
	IDCard     int64 `json:"id_card" db:"id_card"`
	GivesUses  bool  `json:"gives_uses" db:"gives_uses"`
}

// Create a link between an element and a card.
func CreateElementLink(db *gorp.DbMap, card *Card, elem *Element, GivesUses bool) (*ElementLink, error) {
	if db == nil || elem == nil || card == nil {
		return nil, errors.New("Missing parameters to create element link")
	}

	el := &ElementLink{
		IDScenario: card.IDScenario,
		IDElement:  elem.ID,
		IDCard:     card.ID,
		GivesUses:  GivesUses,
	}

	err := db.Insert(el)
	if err != nil {
		return nil, err
	}

	return el, nil
}

// List element links, with filters.
func ListElementLinks(db *gorp.DbMap, scenar *Scenario, card *Card, elem *Element) ([]*ElementLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load element links")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"element_link"`)

	if scenar != nil {
		selector.Where(
			squirrel.Eq{`id_scenario`: scenar.ID},
		)
	}
	if card != nil {
		selector.Where(
			squirrel.Eq{`id_card`: card.ID},
		)
	}
	if elem != nil {
		selector.Where(
			squirrel.Eq{`id_element`: elem.ID},
		)
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var el []*ElementLink

	_, err = db.Select(&el, query, args...)
	if err != nil {
		return nil, err
	}

	return el, nil
}

// Load an element link by id, with optional scenario filter.
func LoadElementLinkFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*ElementLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load element link")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"element_link"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar != nil {
		selector.Where(
			squirrel.Eq{`id_scenario`: scenar.ID},
		)
	}

	query, args, err := selector.ToSql()

	var el ElementLink

	err = db.SelectOne(&el, query, args...)
	if err != nil {
		return nil, err
	}

	return &el, nil
}

// Delete an element link
func (el *ElementLink) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete element link")
	}

	rows, err := db.Delete(el)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such element link to delete")
	}

	return nil
}
