package card

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/models/icon"
	"github.com/loopfz/scecret/models/scenario"
	"github.com/loopfz/scecret/models/stat"
	"github.com/loopfz/scecret/models/statetoken"
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
	TextAreaSize int
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

// SkillTest represents a statistic skill-test required by a Card.
// It links the Card and the Stat for easy querying/navigation, and allows easy
// declaration of test difficulty by declaring each number of shields.
// The code managing SkillTest objects will create CardIcon objects for
// the skill-test statistic and for each kind of non-zero Shield, and link them to the Card.
// That way, the user only has to create a SkillTest to have all the necessary graphical
// elements added to a Card, which can then be edited individually.
type SkillTest struct {
	ID             int64 `json:"id" db:"id"`
	IDCard         int64 `json:"id_card" db:"id_card"`
	IDStat         int64 `json:"id_stat" db:"id_stat"`
	NormalShields  uint  `json:"normal_shields" db:"normal_shields"`
	SkullShields   uint  `json:"skull_shields" db:"skull_shields"`
	HeartShields   uint  `json:"heart_shields" db:"heart_shields"`
	UTShields      uint  `json:"ut_shields" db:"ut_shields"`
	SpecialShields uint  `json:"special_shields" db:"special_shields"`
}

// CardLink represents a Card -> StateToken relation.
// It is uni-directional but can be configured using the UnlocksUnlocked parameter.
// This lets the user express relations such as:
//
// Card X (unlocks) -> StateToken Y (unlocks) -> Card Z
//
// A CardIcon object containing the StateToken icon will be added to the Front of any Card that
// "unlocks" it, and to the back of any Card that is "unlocked" by it.
type CardLink struct {
	ID              int64 `json:"id" db:"id"`
	IDCard          int64 `json:"id_card" db:"id_card"`
	IDStateToken    int64 `json:"id_state_token" db:"id_state_token"`
	UnlocksUnlocked bool  `json:"unlocks_unlocked" db:"unlocks_unlocked"`
}

// CardIcon represents a single graphical icon on the Front or the Back of a Card.
// It is linked to a Card, and to a collection of Icon graphical elements.
// It has coordinates/size properties, and optional annotations (small circle or square) to add
// e.g. a Stat value for a character or a number above a Shield.
// It also has foreign keys to the SkillTest/CardLink that it originated from,
// so that it can be retrieved for update when the SkillTest/CardLink is updated, and can
// be automatically deleted through CASCADE.
type CardIcon struct {
	ID             int64  `json:"id" db:"id"`
	IDCard         int64  `json:"id_card" db:"id_card"`
	FrontBack      bool   `json:"front_back" db:"front_back"`
	IDIcon         int64  `json:"id_icon" db:"id_icon"`
	X              uint   `json:"x" db:"x"`
	Y              uint   `json:"y" db:"y"`
	SizeX          uint   `json:"size_x" db:"size_x"`
	SizeY          uint   `json:"size_y" db:"size_y"`
	Annotation     string `json:"annotation" db:"annotation"`
	AnnotationType int    `json:"annotation_type" db:"annotation_type"`
	IDSkillTest    *int64 `json:"-" db:"id_skilltest"`
	IDCardLink     *int64 `json:"-" db:"id_cardlink"`
}

/*
** BASE CARD
 */

// Create a card.
func Create(db *gorp.DbMap, scenar *scenario.Scenario, num uint, desc string, front *CardFace, back *CardFace) (*Card, error) {
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

// Load a card by ID.
// Scenario parameter is required and acts as a filter.
func LoadFromID(db *gorp.DbMap, scenar *scenario.Scenario, ID int64) (*Card, error) {
	if db == nil || scenar == nil {
		return nil, errors.New("Missing parameters to load card")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"card"`).Where(
		squirrel.And{
			squirrel.Eq{`id_scenario`: scenar.ID},
			squirrel.Eq{`id`: ID},
		},
	).ToSql()
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
** SKILL TEST
 */

// Create a skill test.
// This will also create CardIcon objects on the Front of the Card, for the statistic itself and each of the present shields.
func (c *Card) CreateSkillTest(db *gorp.DbMap, linkedStat *stat.Stat, NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) (*SkillTest, error) {
	if db == nil || linkedStat == nil {
		return nil, errors.New("Missing parameters to create skill test")
	}

	st := &SkillTest{
		IDCard:         c.ID,
		IDStat:         linkedStat.ID,
		NormalShields:  NormalShields,
		SkullShields:   SkullShields,
		HeartShields:   HeartShields,
		SpecialShields: SpecialShields,
	}

	err := st.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(st)
	if err != nil {
		return nil, err
	}

	err = addShieldsAndStatIcon(db, c, linkedStat, st,
		NormalShields, SkullShields, HeartShields, UTShields, SpecialShields)
	if err != nil {
		return nil, err // TODO Tx
	}

	return st, nil
}

func addShieldsAndStatIcon(db *gorp.DbMap, c *Card, linkedStat *stat.Stat, st *SkillTest,
	NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) error {
	// Add shield CardIcons
	offsetX, err := addShieldCardIcon(db, c, 0, NormalShields, icon.NORMAL_SHIELD, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, SkullShields, icon.SKULL_SHIELD, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, HeartShields, icon.HEART_SHIELD, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, HeartShields, icon.UT_SHIELD, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, SpecialShields, icon.SPECIAL_SHIELD, st)
	if err != nil {
		return err
	}
	// Add stat CardIcon
	ico, err := icon.LoadFromID(db, nil, linkedStat.IDIcon)
	if err != nil {
		return err
	}
	_, err = c.CreateCardIcon(db, ico, true, /* FRONT */
		offsetX, 0, DEFAULT_SIZE_X, DEFAULT_SIZE_Y, "", 0, st, nil)
	if err != nil {
		return err
	}
	return nil
}

func addShieldCardIcon(db *gorp.DbMap, c *Card, offsetX uint, shieldCount uint, shieldShortName string, st *SkillTest) (uint, error) {
	if shieldCount == 0 {
		return offsetX, nil
	}

	ico, err := icon.LoadBaseIconFromShortName(db, shieldShortName)
	if err != nil {
		return offsetX, err
	}

	var annot string
	var annotType int
	if shieldCount > 1 {
		annot = strconv.FormatUint(uint64(shieldCount), 10)
		annotType = AnnotationTypeCircle
	}
	_, err = c.CreateCardIcon(db, ico, true, /* FRONT */
		offsetX, 0, DEFAULT_SIZE_X, DEFAULT_SIZE_Y, annot, annotType, st, nil)
	if err != nil {
		return offsetX, err // TODO Tx
	}

	return offsetX + DEFAULT_SIZE_X, nil
}

// List all skill tests linked to a card.
func (c *Card) ListSkillTests(db *gorp.DbMap) ([]*SkillTest, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load skill tests")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"skill_test"`).Where(
		squirrel.Eq{`id_card`: c.ID},
	).ToSql()

	var st []*SkillTest

	_, err = db.Select(&st, query, args...)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// Load a skill test linked to a card, by ID.
func (c *Card) LoadSkillTestByID(db *gorp.DbMap, IDSkillTest int64) (*SkillTest, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load skill test")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"skill_test"`).Where(
		squirrel.And{
			squirrel.Eq{`id`: IDSkillTest},
			squirrel.Eq{`id_card`: c.ID},
		},
	).ToSql()

	var st SkillTest

	err = db.SelectOne(&st, query, args...)
	if err != nil {
		return nil, err
	}
	return &st, nil
}

// Update a skill test linked to a card.
// This will also create CardIcon objects on the Front of the Card, for the statistic itself and each of the present shields.
// The card parameter CANNOT overwrite the card associated with the SkillTest, it is permanent.
// It is passed only to be able to retrieve the associated CardIcon objects.
func (st *SkillTest) Update(db *gorp.DbMap, card *Card, linkedStat *stat.Stat, NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) error {
	if db == nil || linkedStat == nil {
		return errors.New("Missing parameters to update skill test")
	}

	st.IDStat = linkedStat.ID
	st.NormalShields = NormalShields
	st.SkullShields = SkullShields
	st.HeartShields = HeartShields
	st.SpecialShields = SpecialShields

	err := st.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(st)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such skill test to update")
	}

	// Ad-hoc delete of all previous CardIcons
	_, err = db.Exec(`DELETE FROM "card_icon" WHERE id_skilltest = $1`, st.ID)
	if err != nil {
		return err // TODO Tx
	}
	// Recreate new icons
	err = addShieldsAndStatIcon(db, card, linkedStat, st,
		NormalShields, SkullShields, HeartShields, UTShields, SpecialShields)
	if err != nil {
		return err // TODO Tx
	}
	// TODO smarter process to avoid unnecessarily deleting all CardIcons ?

	return nil
}

// Delete a skill test linked to a card.
func (st *SkillTest) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete skill test")
	}

	rows, err := db.Delete(st)
	if err == nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such skill test to delete")
	}
	return nil
}

// Verify that a skill test object is valid before creating/updating it.
func (st *SkillTest) Valid() error {
	if st.IDCard == 0 {
		return errors.New("Missing reference to card object")
	}
	if st.IDStat == 0 {
		return errors.New("Missing reference to stat object")
	}
	if st.NormalShields > MAX_SHIELDS || st.SkullShields > MAX_SHIELDS || st.HeartShields > MAX_SHIELDS || st.SpecialShields > MAX_SHIELDS {
		return fmt.Errorf("Too many shields: max %d", MAX_SHIELDS)
	}
	return nil
}

/*
** CARD LINK
 */

// Create a card link.
// This will also create a CardIcon object representing the state token
// on either the front or the back of the card (depending on if it unlocks / is unlocked).
func (c *Card) CreateCardLink(db *gorp.DbMap, tk *statetoken.StateToken, UnlocksUnlocked bool) (*CardLink, error) {
	if db == nil || tk == nil {
		return nil, errors.New("Missing parameters to create card link")
	}

	cl := &CardLink{
		IDCard:          c.ID,
		IDStateToken:    tk.ID,
		UnlocksUnlocked: UnlocksUnlocked,
	}

	err := db.Insert(cl)
	if err != nil {
		return nil, err
	}

	// Load state token icon
	ico, err := icon.LoadFromID(db, nil, tk.IDIcon)
	if err != nil {
		return nil, err // TODO Tx
	}

	// Create CardIcon of state token icon, on front or back
	_, err = c.CreateCardIcon(db, ico, UnlocksUnlocked, // Unlocks = Front, Unlocked = back
		0, 0, DEFAULT_SIZE_X, DEFAULT_SIZE_Y, "", 0, nil, cl)
	if err != nil {
		return nil, err // TODO Tx
	}

	return cl, nil
}

// List all this card's links.
func (c *Card) ListCardLinks(db *gorp.DbMap) ([]*CardLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card links")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"card_link"`).Where(
		squirrel.Eq{`id_card`: c.ID},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var cl []*CardLink

	_, err = db.Select(&cl, query, args...)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// Load one of this cards links by ID.
func (c *Card) LoadCardLinkByID(db *gorp.DbMap, ID int64) (*CardLink, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card links")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"card_link"`).Where(
		squirrel.And{
			squirrel.Eq{`id`: ID},
			squirrel.Eq{`id_card`: c.ID},
		},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var cl CardLink

	err = db.SelectOne(&cl, query, args...)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

// Update a card link.
// This will also create a CardIcon object representing the state token
// on either the front or the back of the card (depending on if it unlocks / is unlocked).
// The card parameter CANNOT overwrite the card associated with the CardLink, it is permanent.
// It is passed only to be able to retrieve the associated CardIcon object.
func (cl *CardLink) Update(db *gorp.DbMap, card *Card, tk *statetoken.StateToken, UnlocksUnlocked bool) error {
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
		return errors.New("No such card link to update")
	}

	// Retrieve the CardIcons linked to this card + CardLink
	ciList, err := card.ListCardIcons(db, nil, cl)
	if err != nil {
		return err // TODO Tx
	}
	// There should only ever be 1: the state token icon
	if len(ciList) != 1 {
		// Something very wrong happened
		return fmt.Errorf("Invalid state: %d card icons linked to CardLink") // TODO Tx
	}

	ci := ciList[0]
	// Load state token icon
	ico, err := icon.LoadFromID(db, nil, tk.IDIcon)
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

// Delete a card link.
func (cl *CardLink) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete card link")
	}

	rows, err := db.Delete(cl)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such card link to delete")
	}

	return nil
}

/*
** CARD ICON
 */

// Create a CardIcon object.
func (c *Card) CreateCardIcon(db *gorp.DbMap, ico *icon.Icon,
	FrontBack bool, X, Y, SizeX, SizeY uint,
	Annotation string, AnnotationType int,
	SkillTest *SkillTest, CardLink *CardLink) (*CardIcon, error) {
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
	if CardLink != nil {
		ci.IDCardLink = &CardLink.ID
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
func (c *Card) ListCardIcons(db *gorp.DbMap, SkillTest *SkillTest, CardLink *CardLink) ([]*CardIcon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load card icons")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"card_icon"`).Where(
		squirrel.Eq{`id_card`: c.ID},
	)

	if SkillTest != nil {
		selector.Where(squirrel.Eq{`id_skilltest`: SkillTest.ID})
	}
	if CardLink != nil {
		selector.Where(squirrel.Eq{`id_cardlink`: CardLink.ID})
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
func (ci *CardIcon) Update(db *gorp.DbMap, ico *icon.Icon, FrontBack bool,
	X, Y, SizeX, SizeY uint,
	Annotation string, AnnotationType int,
	SkillTest *SkillTest, CardLink *CardLink) error {

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
	if CardLink != nil {
		ci.IDCardLink = &CardLink.ID
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
	if ci.IDSkillTest != nil && ci.IDCardLink != nil {
		return errors.New("References to both skill_test and card_link")
	}
	return nil
}
