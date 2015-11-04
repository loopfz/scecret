package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

const (
	MAX_SHIELDS = 50 // Arbitrarily large

	MAX_X_COORD = 300 // TODO fix
	MAX_Y_COORD = 300 // TODO fix

	DEFAULT_SIZE_X = 20
	DEFAULT_SIZE_Y = 20

	AnnotationTypeSquare = 1
	AnnotationTypeCircle = 2
)

/*
** TYPE DEFINITIONS
 */

// Type card is the basic building block for other game objects.
// It is a generic representation of a card.
type Card struct {
	ID          int64     `json:"id" db:"id"`
	IDScenario  int64     `json:"-" db:"id_scenario"`
	Number      uint      `json:"number" db:"number"`
	Description string    `json:"description" db:"description"`
	Front       *CardFace `json:"front" db:"front"`
	Back        *CardFace `json:"back" db:"back"`
}

// CardFace describes the Front or Back non-image components of a card.
// It is stored as JSON to allow flexibility / evolution in the definition.
type CardFace struct {
	TextAreaSize int         `json:"text_area_size"`
	TextFields   []TextField `json:"text_fields"`
}

// TextField describes a single field of text on a card.
// For now, it just has spacial coordinates.
// Later, this will evolve to contain more flexibility (box_size x/y, font, font_size, opacity, color, ...)
type TextField struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	Text string `json:"text"`
}

// CardIcon represents a single graphical icon on the Front or the Back of a Card.
// It is linked to a Card, and to a collection of Icon graphical elements.
// It has coordinates/size properties, and optional annotations (small circle or square) to add
// e.g. a Stat value for a character or a number above a Shield.
// It also has foreign keys to the SkillTest/StateTokenLink that it originated from,
// so that it can be retrieved for update when the SkillTest/StateTokenLink is updated, and can
// be automatically deleted through CASCADE.
type CardIcon struct {
	ID               int64  `json:"id" db:"id"`
	IDCard           int64  `json:"id_card" db:"id_card"`
	FrontBack        bool   `json:"front_back" db:"front_back"`
	IDIcon           int64  `json:"id_icon" db:"id_icon"`
	X                uint   `json:"x" db:"x"`
	Y                uint   `json:"y" db:"y"`
	SizeX            uint   `json:"size_x" db:"size_x"`
	SizeY            uint   `json:"size_y" db:"size_y"`
	Annotation       string `json:"annotation" db:"annotation"`
	AnnotationType   int    `json:"annotation_type" db:"annotation_type"`
	IDSkillTest      *int64 `json:"-" db:"id_skilltest"`
	IDStateTokenLink *int64 `json:"-" db:"id_statetokenlink"`
}

/*
** BASE CARD
 */

// Create a card.
func CreateCard(db *gorp.DbMap, scenar *Scenario, num uint, desc string, front *CardFace, back *CardFace) (*Card, error) {
	if db == nil || scenar == nil {
		return nil, errors.New("Missing parameters to create card")
	}

	c := &Card{
		IDScenario:  scenar.ID,
		Number:      num,
		Description: desc,
		Front:       front,
		Back:        back,
	}

	err := db.Insert(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Load a card by ID. Optionally filtered by scenario.
func LoadCardFromID(db *gorp.DbMap, scenar *Scenario, ID int64) (*Card, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"card"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar != nil {
		selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var c Card

	err = db.SelectOne(&c, query, args...)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// Update a card.
func (c *Card) Update(db *gorp.DbMap, num uint, desc string, front *CardFace, back *CardFace) error {
	if db == nil {
		return errors.New("Missing db parameter to update card")
	}

	c.Number = num
	c.Description = desc
	c.Front = front
	c.Back = back

	rows, err := db.Update(c)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such card to update")
	}
	return nil
}

// Delete a card.
func (c *Card) Delete(db *gorp.DbMap) error {
	rows, err := db.Delete(c)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such card to delete")
	}
	return nil
}

/*
** CARD ICON
 */

// Create a CardIcon object.
func (c *Card) CreateCardIcon(db *gorp.DbMap, ico *Icon,
	FrontBack bool, X, Y, SizeX, SizeY uint,
	Annotation string, AnnotationType int,
	SkillTest *SkillTest, StateTokenLink *StateTokenLink) (*CardIcon, error) {
	if db == nil || ico == nil {
		return nil, errors.New("Missing parameters to create card icon")
	}

	ci := &CardIcon{
		IDCard:         c.ID,
		FrontBack:      FrontBack,
		IDIcon:         ico.ID,
		X:              X,
		Y:              Y,
		SizeX:          SizeX,
		SizeY:          SizeY,
		Annotation:     Annotation,
		AnnotationType: AnnotationType,
	}

	if SkillTest != nil {
		ci.IDSkillTest = &SkillTest.ID
	}
	if StateTokenLink != nil {
		ci.IDStateTokenLink = &StateTokenLink.ID
	}

	err := ci.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(ci)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

// List all CardIcon objects linked to this card, with filters.
func (c *Card) ListCardIcons(db *gorp.DbMap, SkillTest *SkillTest, StateTokenLink *StateTokenLink) ([]*CardIcon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card icons")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"card_icon"`).Where(
		squirrel.Eq{`id_card`: c.ID},
	)

	if SkillTest != nil {
		selector.Where(squirrel.Eq{`id_skilltest`: SkillTest.ID})
	}
	if StateTokenLink != nil {
		selector.Where(squirrel.Eq{`id_statetokenlink`: StateTokenLink.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var ci []*CardIcon

	_, err = db.Select(&ci, query, args...)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

// Load one CardIcon object linked to this card, by ID.
func (c *Card) LoadCardIconByID(db *gorp.DbMap, ID int64) (*CardIcon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card icon")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"card_icon"`).Where(
		squirrel.And{
			squirrel.Eq{`id`: ID},
			squirrel.Eq{`id_card`: c.ID},
		},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var ci CardIcon

	err = db.SelectOne(&ci, query, args...)
	if err != nil {
		return nil, err
	}

	return &ci, nil
}

// Update a CardIcon object.
func (ci *CardIcon) Update(db *gorp.DbMap, ico *Icon, FrontBack bool,
	X, Y, SizeX, SizeY uint,
	Annotation string, AnnotationType int,
	SkillTest *SkillTest, StateTokenLink *StateTokenLink) error {

	ci.FrontBack = FrontBack
	ci.IDIcon = ico.ID
	ci.X = X
	ci.Y = Y
	ci.SizeX = SizeX
	ci.SizeY = SizeY
	ci.Annotation = Annotation
	ci.AnnotationType = AnnotationType
	if SkillTest != nil {
		ci.IDSkillTest = &SkillTest.ID
	}
	if StateTokenLink != nil {
		ci.IDStateTokenLink = &StateTokenLink.ID
	}

	err := ci.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(ci)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such card icon to update")
	}

	return nil
}

// Delete a CardIcon object.
func (ci *CardIcon) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete card icon")
	}

	rows, err := db.Delete(ci)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such card icon to delete")
	}

	return nil
}

// Verify that a CardIcon object is valid before creating/updating it.
func (ci *CardIcon) Valid() error {
	if ci.IDCard == 0 {
		return errors.New("Missing reference to card object")
	}
	if ci.IDIcon == 0 {
		return errors.New("Missing reference to icon object")
	}
	if ci.X+ci.SizeX > MAX_X_COORD {
		return fmt.Errorf("X coord: %d too big (max %d)", ci.X+ci.SizeX, MAX_X_COORD)
	}
	if ci.Y+ci.SizeY > MAX_Y_COORD {
		return fmt.Errorf("Y coord: %d too big (max %d)", ci.Y+ci.SizeY, MAX_Y_COORD)
	}
	if ci.AnnotationType != 0 && ci.AnnotationType != AnnotationTypeSquare && ci.AnnotationType != AnnotationTypeCircle {
		return fmt.Errorf("Unknown annotation type %d", ci.AnnotationType)
	}
	if ci.IDSkillTest != nil && ci.IDStateTokenLink != nil {
		return errors.New("References to both skill_test and state_token_link")
	}
	return nil
}

func (cf *CardFace) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	s := value.([]byte)
	return json.Unmarshal(s, &cf)
}

func (cf *CardFace) Value() (driver.Value, error) {
	if cf == nil {
		return nil, nil
	}
	j, err := json.Marshal(cf)
	if err != nil {
		return nil, err
	}
	return j, nil
}
