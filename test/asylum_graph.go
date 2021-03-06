package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/loopfz/scecret/db/testdb"
	"github.com/loopfz/scecret/models"
)

func main() {

	db, err := testdb.InitTestDB(true)
	if err != nil {
		panic(err)
	}

	scenar, err := models.CreateScenario(db, "Asylum", &models.User{})
	if err != nil {
		panic(err)
	}

	elemDortoirMap, err := models.CreateElement(db, scenar, 1, "Plan vers le cabinet")
	if err != nil {
		panic(err)
	}
	elemDortoirMapCard, err := models.LoadCardFromID(db, scenar, elemDortoirMap.IDCard)
	if err != nil {
		panic(err)
	}

	// Ad-hoc creation, missing model code
	combatIcon := &models.Icon{IDScenario: &scenar.ID}
	err = db.Insert(combatIcon)
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
	_, dortoirCards := createLoc(db, scenar, "Dortoir", 6)

	// CABINET
	cabinet, cabinetCards := createLoc(db, scenar, "Cabinet", 5)

	// LABYRINTHE
	labyrinthe, labyrintheCards := createLoc(db, scenar, "Labyrinthe", 5)

	// PARC
	parc, parcCards := createLoc(db, scenar, "Parc", 5)

	// SERRE
	serre, serreCards := createLoc(db, scenar, "Serre", 4)

	// TOMBEAU
	tombeau, tombeauCards := createLoc(db, scenar, "Tombeau", 5)

	// PORTE PENTACLES
	portePentacles, portePentaclesCards := createLoc(db, scenar, "Porte pentacles", 4)

	// CATACOMBES
	catacombes, catacombesCards := createLoc(db, scenar, "Catacombes", 5)

	// CRYPTE
	crypte, crypteCards := createLoc(db, scenar, "Crypte", 7)

	_, err = models.CreateElementLink(db, dortoirCards[5], elemDortoirMap, true)
	if err != nil {
		panic(err)
	}

	createLocLink(db, infirmerieCards[2], cabinet)
	createLocLink(db, promenadeCards[2], parc)
	//createLocLink(db, dortoirCards[5], catacombes)
	createLocLink(db, elemDortoirMapCard, catacombes)
	createLocLink(db, cabinetCards[3], parc)
	createLocLink(db, cabinetCards[4], labyrinthe)
	createLocLink(db, labyrintheCards[0], parc)
	createLocLink(db, parcCards[4], portePentacles)
	createLocLink(db, parcCards[3], serre)
	createLocLink(db, catacombesCards[4], portePentacles)
	createLocLink(db, portePentaclesCards[1], tombeau)
	createLocLink(db, portePentaclesCards[2], crypte)
	createLocLink(db, portePentaclesCards[3], catacombes)

	createSkillTest(db, cabinetCards[2], combat)
	createSkillTest(db, dortoirCards[4], combat)
	createSkillTest(db, parcCards[1], combat)
	createSkillTest(db, parcCards[4], combat)
	createSkillTest(db, serreCards[3], combat)
	createSkillTest(db, catacombesCards[1], combat)
	createSkillTest(db, catacombesCards[3], combat)
	createSkillTest(db, catacombesCards[4], combat)
	createSkillTest(db, tombeauCards[1], combat)
	createSkillTest(db, tombeauCards[2], combat)
	createSkillTest(db, tombeauCards[3], combat)
	createSkillTest(db, crypteCards[5], combat)
	createSkillTest(db, crypteCards[6], combat)

	createStateTokenLink(db, infirmerieCards[1], stateToken, true)
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

func createSkillTest(db *gorp.DbMap, card *models.Card, stat *models.Stat) {
	_, err := models.CreateSkillTest(db, card, stat, 0, 0, 0, 0, 0)
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
