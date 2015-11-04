package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/models"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	sqldb, err := sql.Open("sqlite3", fmt.Sprintf("/tmp/scecret%d.db", rand.Int()))
	if err != nil {
		panic(err)
	}

	db := &gorp.DbMap{Db: sqldb, Dialect: gorp.SqliteDialect{}}

	db.AddTableWithName(models.Scenario{}, `scenario`).SetKeys(true, "id")
	db.AddTableWithName(models.Location{}, `location`).SetKeys(true, "id")
	db.AddTableWithName(models.LocationCard{}, `location_card`).SetKeys(true, "id")
	db.AddTableWithName(models.LocationLink{}, `location_link`).SetKeys(true, "id")
	db.AddTableWithName(models.Card{}, `card`).SetKeys(true, "id")
	db.AddTableWithName(models.CardIcon{}, `card_icon`).SetKeys(true, "id")
	db.AddTableWithName(models.Element{}, `element`).SetKeys(true, "id")
	db.AddTableWithName(models.Icon{}, `icon`).SetKeys(true, "id")
	db.AddTableWithName(models.StateToken{}, `state_token`).SetKeys(true, "id")
	db.AddTableWithName(models.StateTokenLink{}, `state_token_link`).SetKeys(true, "id")
	db.AddTableWithName(models.Stat{}, `stat`).SetKeys(true, "id")
	db.AddTableWithName(models.SkillTest{}, `skill_test`).SetKeys(true, "id")

	err = db.CreateTablesIfNotExists()
	if err != nil {
		panic(err)
	}

	scenar, err := models.CreateScenario(db, "Asylum", &models.User{})
	if err != nil {
		panic(err)
	}

	// Ad-hoc creation, missing model code
	combatIcon := &models.Icon{IDScenario: &scenar.ID}
	err = db.Insert(combatIcon)
	if err != nil {
		panic(err)
	}
	blockingIcon := &models.Icon{ShortName: models.BLOCKING_ICON}
	err = db.Insert(blockingIcon)
	if err != nil {
		panic(err)
	}
	stateTokenIcon := &models.Icon{ShortName: "STATE_TOKEN_TEST"}
	err = db.Insert(stateTokenIcon)
	if err != nil {
		panic(err)
	}
	combat := &models.Stat{Name: "Combat", IDScenario: scenar.ID, IDIcon: combatIcon.ID}
	err = db.Insert(combat)
	if err != nil {
		panic(err)
	}
	stateToken := &models.StateToken{
		ShortName: "STATE_TOKEN_TEST",
		IDIcon:    stateTokenIcon.ID,
	}
	err = db.Insert(stateToken)
	if err != nil {
		panic(err)
	}

	// REPOS
	createLoc(db, scenar, "Repos", 6)

	// INFIRMERIE
	_, infirmerieCards := createLoc(db, scenar, "Infirmerie", 5)

	// PROMENADE
	_, promenadeCards := createLoc(db, scenar, "Promenade", 7)

	// CUISINE
	createLoc(db, scenar, "Cuisine", 4)

	// DORTOIR
	_, dortoirCards := createLoc(db, scenar, "Dortoir", 4)

	// CABINET
	cabinet, cabinetCards := createLoc(db, scenar, "Cabinet", 6)

	// LABYRINTHE
	labyrinthe, labyrintheCards := createLoc(db, scenar, "Labyrinthe", 6)

	// PARC
	parc, parcCards := createLoc(db, scenar, "Parc", 5)

	// SERRE
	serre, serreCards := createLoc(db, scenar, "Serre", 3)

	// TOMBEAU
	tombeau, tombeauCards := createLoc(db, scenar, "Tombeau", 5)

	// PORTE PENTACLES
	portePentacles, portePentaclesCards := createLoc(db, scenar, "Porte pentacles", 4)

	// CATACOMBES
	catacombes, catacombesCards := createLoc(db, scenar, "Catacombes", 5)

	// CRYPTE
	crypte, crypteCards := createLoc(db, scenar, "Crypte", 8)

	createLocLink(db, infirmerieCards[3], cabinet)
	createLocLink(db, promenadeCards[3], parc)
	createLocLink(db, dortoirCards[1], catacombes)
	createLocLink(db, cabinetCards[4], labyrinthe)
	createLocLink(db, cabinetCards[5], parc)
	createLocLink(db, labyrintheCards[5], parc)
	createLocLink(db, parcCards[2], portePentacles)
	createLocLink(db, parcCards[4], serre)
	createLocLink(db, catacombesCards[4], portePentacles)
	createLocLink(db, portePentaclesCards[1], catacombes)
	createLocLink(db, portePentaclesCards[2], crypte)
	createLocLink(db, portePentaclesCards[3], tombeau)

	createSkillTest(db, cabinetCards[4], combat, true)
	createSkillTest(db, parcCards[1], combat, true)
	createSkillTest(db, serreCards[2], combat, true)
	createSkillTest(db, catacombesCards[1], combat, true)
	createSkillTest(db, catacombesCards[3], combat, true)
	createSkillTest(db, catacombesCards[4], combat, true)
	createSkillTest(db, tombeauCards[1], combat, true)
	createSkillTest(db, tombeauCards[2], combat, true)
	createSkillTest(db, tombeauCards[3], combat, true)
	createSkillTest(db, tombeauCards[4], combat, true)
	createSkillTest(db, crypteCards[4], combat, true)
	createSkillTest(db, crypteCards[5], combat, true)

	createStateTokenLink(db, infirmerieCards[4], stateToken, true)
	createStateTokenLink(db, dortoirCards[3], stateToken, false)

	out, err := models.Graph(db, scenar)
	if err != nil {
		panic(err)
	}

	jsonOut, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonOut))
}

func createLoc(db *gorp.DbMap, scenar *models.Scenario, name string, numCards int) (*models.Location, []*models.Card) {
	loc, err := models.CreateLocation(db, scenar, name, false)
	if err != nil {
		panic(err)
	}
	letters := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	for i, l := range letters {
		if i < numCards {
			createLocCard(db, loc, scenar, l)
		}
	}
	locCards, err := loc.GetCards(db)
	if err != nil {
		panic(err)
	}
	return loc, locCards
}

func createLocCard(db *gorp.DbMap, loc *models.Location, scenar *models.Scenario, letter string) {
	_, err := loc.CreateLocationCard(db, scenar, letter)
	if err != nil {
		panic(err)
	}
}

func createLocLink(db *gorp.DbMap, card *models.Card, loc *models.Location) {
	_, err := models.CreateLocationLink(db, card, loc)
	if err != nil {
		panic(err)
	}
}

func createSkillTest(db *gorp.DbMap, card *models.Card, stat *models.Stat, blocking bool) {
	_, err := models.CreateSkillTest(db, card, stat, blocking, 0, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
}

func createStateTokenLink(db *gorp.DbMap, card *models.Card, st *models.StateToken, unlocksUnlocked bool) {
	_, err := models.CreateStateTokenLink(db, card, st, unlocksUnlocked)
	if err != nil {
		panic(err)
	}
}
