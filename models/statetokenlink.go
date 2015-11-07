package models

import (
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

// StateTokenLink represents a Card -> StateToken relation.
// It is uni-directional but can be configured using the UnlocksUnlocked parameter.
// This lets the user express relations such as:
//
// Card X (unlocks) -> StateToken Y (unlocks) -> Card Z
//
// A CardIcon object containing the StateToken icon will be added to the Front of any Card that
// "unlocks" it, and to the back of any Card that is "unlocked" by it.
type StateTokenLink struct {
	ID              int64 `json:"id" db:"id"`
	IDScenario      int64 `json:"-" db:"id_scenario"`
	IDCard          int64 `json:"id_card" db:"id_card"`
	IDStateToken    int64 `json:"id_state_token" db:"id_state_token"`
	UnlocksUnlocked bool  `json:"unlocks_unlocked" db:"unlocks_unlocked"`
}

// Create a link between a card and a state token.
// This will also create a CardIcon object representing the state token
// on either the front or the back of the card (depending on if it unlocks / is unlocked).
func CreateStateTokenLink(db *gorp.DbMap, card *Card, tk *StateToken, UnlocksUnlocked bool) (*StateTokenLink, error) {
	if db == nil || tk == nil || card == nil {
		return nil, errors.New("Missing parameters to create card link")
	}

	cl := &StateTokenLink{
		IDScenario:      card.IDScenario,
		IDCard:          card.ID,
		IDStateToken:    tk.ID,
		UnlocksUnlocked: UnlocksUnlocked,
	}

	err := db.Insert(cl)
	if err != nil {
		return nil, err
	}

	// Load state token icon
	ico, err := LoadIconFromID(db, nil, tk.IDIcon)
	if err != nil {
		return nil, err // TODO Tx
	}

	// Create CardIcon of state token icon, on front or back
	_, err = card.CreateCardIcon(db, ico, UnlocksUnlocked, // Unlocks = Front, Unlocked = back
		0, 0, DEFAULT_SIZE_X, DEFAULT_SIZE_Y, "", 0, nil, cl)
	if err != nil {
		return nil, err // TODO Tx
	}

	return cl, nil
}

// List state token links, with filters.
func ListStateTokenLinks(db *gorp.DbMap, scenar *Scenario, card *Card) ([]*StateTokenLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load state token links")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"state_token_link"`)

	if scenar != nil {
		selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}
	if card != nil {
		selector.Where(squirrel.Eq{`id_card`: card.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var cl []*StateTokenLink

	_, err = db.Select(&cl, query, args...)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// Loads a state token link by ID. Optionally filtered by scenario.
func LoadStateTokenLinkByID(db *gorp.DbMap, scenar *Scenario, ID int64) (*StateTokenLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card links")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"state_token_link"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar != nil {
		selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var cl StateTokenLink

	err = db.SelectOne(&cl, query, args...)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

// TODO delete?
// Update a card link.
// This will also create a CardIcon object representing the state token
// on either the front or the back of the card (depending on if it unlocks / is unlocked).
// The card parameter CANNOT overwrite the card associated with the StateTokenLink, it is permanent.
// It is passed only to be able to retrieve the associated CardIcon object.
func (cl *StateTokenLink) Update(db *gorp.DbMap, card *Card, tk *StateToken, UnlocksUnlocked bool) error {
	if db == nil || tk == nil {
		return errors.New("Missing parameters to create card link")
	}

	cl.IDStateToken = tk.ID
	cl.UnlocksUnlocked = UnlocksUnlocked

	rows, err := db.Update(cl)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such state token link to update")
	}

	// Retrieve the CardIcons linked to this card + StateTokenLink
	ciList, err := card.ListCardIcons(db, nil, cl)
	if err != nil {
		return err // TODO Tx
	}
	// There should only ever be 1: the state token icon
	if len(ciList) != 1 {
		// Something very wrong happened
		return fmt.Errorf("Invalid state: %d card icons linked to StateTokenLink") // TODO Tx
	}

	ci := ciList[0]
	// Load state token icon
	ico, err := LoadIconFromID(db, nil, tk.IDIcon)
	if err != nil {
		return err // TODO Tx
	}

	// Update the CardIcon, changing only the ico parameter
	err = ci.Update(db, ico, ci.FrontBack, ci.X, ci.Y, ci.SizeX, ci.SizeY, ci.Annotation,
		ci.AnnotationType, nil, cl)
	if err != nil {
		return err // TODO Tx
	}

	return nil
}

// Delete a state token link.
func (cl *StateTokenLink) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete state token link")
	}

	rows, err := db.Delete(cl)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such state token link to delete")
	}

	return nil
}
