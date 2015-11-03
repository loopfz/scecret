package icon

import (
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/models/scenario"
	"github.com/loopfz/scecret/utils/sqlgenerator"
)

const (
	// Const names to be retrievable from model code
	NORMAL_SHIELD  = "normal_shield"
	SKULL_SHIELD   = "skull_shield"
	HEART_SHIELD   = "heart_shield"
	UT_SHIELD      = "ut_shield"
	SPECIAL_SHIELD = "special_shield"
)

type Icon struct {
	ID         int64  `json:"id" db:"id"`
	IDScenario *int64 `json:"id_scenario" db:"id_scenario"`
	ShortName  string `json:"short_name" db:"short_name"`
	URL        string `json:"url" db:"url"`
}

// Create an icon
func Create(db *gorp.DbMap, scenar *scenario.Scenario, ShortName string, URL string) (*Icon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to create icon")
	}

	i := &Icon{
		ShortName: ShortName,
		URL:       URL,
	}

	if scenar != nil {
		i.IDScenario = &scenar.ID
	}

	err := i.Valid()
	if err != nil {
		return nil, err
	}

	err = db.Insert(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

// Load an icon from ID. If scenar parameter is non-nil it acts as a filter: only rows with id_scenario NULL or stricly equal
// will be returned.
func LoadFromID(db *gorp.DbMap, scenar *scenario.Scenario, ID int64) (*Icon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load icon")
	}

	selector := sqlgenerator.PGsql.Select(`*`).From(`"icon"`).Where(
		squirrel.Eq{`id`: ID},
	)

	if scenar == nil {
		selector.Where(squirrel.Eq{`id_scenario`: nil}) // squirrel.Eq properly transforms nil rvalues into IS NULL
	} else {
		selector.Where(squirrel.Or{squirrel.Eq{`id_scenario`: nil}, squirrel.Eq{`id_scenario`: scenar.ID}})
	}

	query, args, err := selector.ToSql()
	if err != nil {
		return nil, err
	}

	var i Icon

	err = db.SelectOne(&i, query, args...)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// Used to load base game objects, e.g. shield icons. These need to be referenced by a const name for conveniency.
// This enforces id_scenario IS NULL (i.e. base game objects) on returned rows.
func LoadBaseIconFromShortName(db *gorp.DbMap, ShortName string) (*Icon, error) {
	if db == nil {
		return nil, errors.New("Missing db parameter to load base icon")
	}

	query, args, err := sqlgenerator.PGsql.Select(`*`).From(`"icon"`).Where(
		squirrel.And{
			squirrel.Eq{`id_scenario`: nil}, // squirrel.Eq properly transforms nil rvalues into IS NULL
			squirrel.Eq{`short_name`: ShortName},
		},
	).ToSql()

	if err != nil {
		return nil, err
	}

	var i Icon

	err = db.SelectOne(&i, query, args...)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

// Update an icon
func (i *Icon) Update(db *gorp.DbMap, ShortName string, URL string) error {
	if db == nil {
		return errors.New("Missing db parameter to update icon")
	}

	i.ShortName = ShortName
	i.URL = URL

	err := i.Valid()
	if err != nil {
		return err
	}

	rows, err := db.Update(i)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such icon to update")
	}

	return nil
}

// Delete an icon
func (i *Icon) Delete(db *gorp.DbMap) error {
	if db == nil {
		return errors.New("Missing db parameter to delete icon")
	}

	rows, err := db.Delete(i)
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("No such icon to delete")
	}

	return nil
}

// Verify that an icon object is valid before creating/updating it
func (i *Icon) Valid() error {
	// TODO coherency checks
	return nil
}
