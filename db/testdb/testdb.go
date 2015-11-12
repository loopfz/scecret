package testdb

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/models"
	_ "github.com/mattn/go-sqlite3"
)

func InitTestDB(random bool) (*gorp.DbMap, error) {

	var s string
	if random {
		rand.Seed(time.Now().UTC().UnixNano())
		s = fmt.Sprintf("/tmp/scecret%d.db", rand.Int())
	} else {
		s = "/tmp/scecret.db"
	}
	sqldb, err := sql.Open("sqlite3", s)
	if err != nil {
		return nil, err
	}

	db := &gorp.DbMap{Db: sqldb, Dialect: gorp.SqliteDialect{}}

	db.AddTableWithName(models.User{}, `user`).SetKeys(true, "id")
	db.AddTableWithName(models.Scenario{}, `scenario`).SetKeys(true, "id")
	db.AddTableWithName(models.Location{}, `location`).SetKeys(true, "id")
	db.AddTableWithName(models.LocationCard{}, `location_card`).SetKeys(true, "id")
	db.AddTableWithName(models.LocationLink{}, `location_link`).SetKeys(true, "id")
	db.AddTableWithName(models.Card{}, `card`).SetKeys(true, "id")
	db.AddTableWithName(models.CardIcon{}, `card_icon`).SetKeys(true, "id")
	db.AddTableWithName(models.Element{}, `element`).SetKeys(true, "id")
	db.AddTableWithName(models.ElementLink{}, `element_link`).SetKeys(true, "id")
	db.AddTableWithName(models.Icon{}, `icon`).SetKeys(true, "id")
	db.AddTableWithName(models.StateToken{}, `state_token`).SetKeys(true, "id")
	db.AddTableWithName(models.StateTokenLink{}, `state_token_link`).SetKeys(true, "id")
	db.AddTableWithName(models.Stat{}, `stat`).SetKeys(true, "id")
	db.AddTableWithName(models.SkillTest{}, `skill_test`).SetKeys(true, "id")

	err = db.CreateTablesIfNotExists()
	if err != nil {
		return nil, err
	}

	return db, nil
}
