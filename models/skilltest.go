package models

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

// SkillTest represents a statistic skill-test required by a Card.
// It links the Card and the Stat for easy querying/navigation, and allows easy
// declaration of test difficulty by declaring each number of shields.
// The code managing SkillTest objects will create CardIcon objects for
// the skill-test statistic and for each kind of non-zero Shield, and link them to the Card.
// That way, the user only has to create a SkillTest to have all the necessary graphical
// elements added to a Card, which can then be edited individually.
type SkillTest struct {
	ID             int64 `json:"id" db:"id"`
	IDScenario     int64 `json:"-" db:"id_scenario"`
	IDCard         int64 `json:"id_card" db:"id_card"`
	IDStat         int64 `json:"id_stat" db:"id_stat"`
	NormalShields  uint  `json:"normal_shields" db:"normal_shields"`
	SkullShields   uint  `json:"skull_shields" db:"skull_shields"`
	HeartShields   uint  `json:"heart_shields" db:"heart_shields"`
	UTShields      uint  `json:"ut_shields" db:"ut_shields"`
	SpecialShields uint  `json:"special_shields" db:"special_shields"`
}

// Create a skill test.
// This will also create CardIcon objects on the Front of the Card, for the statistic itself and each of the present shields.
func CreateSkillTest(db *gorp.DbMap, card *Card, linkedStat *Stat, NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) (*SkillTest, error) {
	if db == nil || linkedStat == nil {
		return nil, errors.New("Missing parameters to create skill test")
	}

	st := &SkillTest{
		IDScenario:     card.IDScenario,
		IDCard:         card.ID,
		IDStat:         linkedStat.ID,
		NormalShields:  NormalShields,
		SkullShields:   SkullShields,
		HeartShields:   HeartShields,
		UTShields:      UTShields,
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

	err = addSkillTestIcons(db, card, linkedStat, st,
		NormalShields, SkullShields, HeartShields, UTShields, SpecialShields)
	if err != nil {
		return nil, err // TODO Tx
	}

	return st, nil
}

func addSkillTestIcons(db *gorp.DbMap, c *Card, linkedStat *Stat, st *SkillTest,
	NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) error {
	// Add shield CardIcons
	offsetX, err := addShieldCardIcon(db, c, 0, NormalShields, NORMAL_SHIELD_ICON, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, SkullShields, SKULL_SHIELD_ICON, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, HeartShields, HEART_SHIELD_ICON, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, HeartShields, UT_SHIELD_ICON, st)
	if err != nil {
		return err
	}
	offsetX, err = addShieldCardIcon(db, c, offsetX, SpecialShields, SPECIAL_SHIELD_ICON, st)
	if err != nil {
		return err
	}

	// Add stat CardIcon
	ico, err := LoadIconFromID(db, nil, linkedStat.IDIcon)
	if err != nil {
		return err
	}
	_, err = c.CreateCardIcon(db, ico, true, /* FRONT */
		offsetX, 0, DEFAULT_SIZE_X, DEFAULT_SIZE_Y, "", 0, st, nil)
	if err != nil {
		return err
	}
	offsetX += DEFAULT_SIZE_X

	return nil
}

func addShieldCardIcon(db *gorp.DbMap, c *Card, offsetX uint, shieldCount uint, shieldShortName string, st *SkillTest) (uint, error) {
	if shieldCount == 0 {
		return offsetX, nil
	}

	ico, err := LoadBaseIconFromShortName(db, shieldShortName)
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

// List skill tests with filters.
func ListSkillTests(db *gorp.DbMap, scenar *Scenario, card *Card, s *Stat) ([]*SkillTest, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load skill tests")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"skill_test"`)

	if scenar != nil {
		selector = selector.Where(
			squirrel.Eq{`id_scenario`: scenar.ID},
		)
	}
	if card != nil {
		selector = selector.Where(squirrel.Eq{`id_card`: card.ID})
	}
	if s != nil {
		selector = selector.Where(squirrel.Eq{`id_stat`: s.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var st []*SkillTest

	_, err = db.Select(&st, query, args...)
	if err != nil {
		return nil, err
	}
	return st, nil
}

// Load a skill test, by ID. Optionally filtered by scenario.
func LoadSkillTestFromID(db *gorp.DbMap, scenar *Scenario, IDSkillTest int64) (*SkillTest, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load skill test")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"skill_test"`).Where(
		squirrel.Eq{`id`: IDSkillTest},
	)

	if scenar != nil {
		selector = selector.Where(squirrel.Eq{`id_scenario`: scenar.ID})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

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
func (st *SkillTest) Update(db *gorp.DbMap, card *Card, linkedStat *Stat, NormalShields, SkullShields, HeartShields, UTShields, SpecialShields uint) error {
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
	err = addSkillTestIcons(db, card, linkedStat, st,
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
