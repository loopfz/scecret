package models

import (
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

type Element struct {
	ID          int64  `json:"-" db:"id"`
	IDScenario  int64  `json:"-" db:"id_scenario"`
	Number      int    `json:"number" db:"number"`
	Description string `json:"description" db:"description"`
	Notes       string `json:"notes" db:"notes"`
	IDCard      int64  `json:"id_card" db:"id_card"`
}

// Create a new element.
func CreateElement(db *gorp.DbMap, scenar *Scenario, Number int, Description string) (*Element, error) {
	if db == nil || scenar == nil {
		return nil, errors.New("Missing parameters to create element")
	}

	cardDesc := fmt.Sprintf("Element %d", Number)
	card, err := CreateCard(db, scenar, 0, cardDesc, &CardFace{}, &CardFace{})
	if err != nil {
		return nil, err
	}

	elem := &Element{
		IDScenario:  scenar.ID,
		Number:      Number,
		Description: Description,
		IDCard:      card.ID,
	}

	err = elem.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(elem)
	if err != nil {
		return nil, err // TODO Tx
	}

	return elem, nil
}

// List elements, optionally filtered by scenario.
func ListElements(db *gorp.DbMap, scenar *Scenario) ([]*Element, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list elements")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"element"`)

	if scenar == nil {
		selector = selector.Where(
			squirrel.Eq{`id_scenario`: scenar.ID},
		)
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var elem []*Element

	_, err = db.Select(&elem, query, args...)
	if err != nil {
		return nil, err
	}

	return elem, nil
}

// Returns all card objects that belong to elements.
func GetElementCards(db *gorp.DbMap, scenar *Scenario) ([]*Card, error) {
	if db == nil || scenar == nil {
		return nil, errors.New("Missing parameters to get element cards")
	}

	selector := sqlgenerator.PGsql.Select(`"card".*`).From(`"card"`).Join(
		`"element" ON "card".id = "element".id_card`,
	).Where(
		squirrel.Eq{`"card".id_scenario`: scenar.ID},
	)

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var c []*Card

	_, err = db.Select(&c, query, args...)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Load element by ID. Optional scenario filter.
func LoadElementFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*Element, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to list elements")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"element"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar == nil {
		selector = selector.Where(
			squirrel.Eq{`id_scenario`: scenar.ID},
		)
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var elem Element

	err = db.SelectOne(&elem, query, args...)
	if err != nil {
		return nil, err
	}

	return &elem, nil
}

// Update an element.
func (e *Element) Update(db *gorp.DbMap, Number int, Description string, Notes string) error {
	if db == nil {
		return errors.New("Missing db parameter to update element")
	}

	e.Number = Number
	e.Description = Description
	e.Notes = Notes

	// Update linked card description
	cardDesc := fmt.Sprintf("Element %d", Number)
	card, err := LoadCardFromID(db, nil, e.IDCard)
	if err != nil {
		return err
	}
	err = card.Update(db, card.Number, cardDesc, card.Front, card.Back)
	if err != nil {
		return err
	}

	err = e.Valid()
	if err != nil {
		return err // TODO Tx
	}

	rows, err := db.Update(e)
	if err != nil {
		return err // TODO Tx
	}
	if rows == 0 {
		return errors.New("No such element to update") // TODO Tx
	}

	return nil
}

// Delete an element.
func (e *Element) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete element")
	}

	// Delete card object
	card, err := LoadCardFromID(db, nil, e.IDCard)
	if err != nil {
		return err
	}
	err = card.Delete(db)
	if err != nil {
		return err
	}

	rows, err := db.Delete(e)
	if err != nil {
		return err // TODO Tx
	}
	if rows == 0 {
		return errors.New("No such element to delete") // TODO Tx
	}

	return nil
}

func (e *Element) Valid() error {
	if e.Number == 0 {
		return errors.New("Missing element number")
	}
	return nil
}
