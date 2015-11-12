package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

func NewSandbox(c *gin.Context) (*models.Scenario, error) {

	u, err := auth.RetrieveTokenUser(db, c)
	if err != nil {
		return nil, err
	}

	sc, err := models.CreateScenario(db, "Asylum sandbox", u)
	if err != nil {
		return nil, err
	}

	locs := map[string]*struct {
		NumA     int
		NumOther int
		Loc      *models.Location
		Cards    []*models.Card
	}{
		"Salle de repos":  {NumA: 1, NumOther: 5},
		"Infirmerie":      {NumA: 1, NumOther: 4},
		"Promenade":       {NumA: 1, NumOther: 6},
		"Cuisine":         {NumA: 1, NumOther: 3},
		"Dortoir":         {NumA: 1, NumOther: 5},
		"Cabinet":         {NumA: 1, NumOther: 4},
		"Labyrinthe":      {NumA: 1, NumOther: 4},
		"Parc":            {NumA: 1, NumOther: 4},
		"Serre":           {NumA: 1, NumOther: 3},
		"Tombeau":         {NumA: 1, NumOther: 4},
		"Porte pentacles": {NumA: 3, NumOther: 3},
		"Catacombes":      {NumA: 1, NumOther: 4},
		"Crypte":          {NumA: 1, NumOther: 4},
	}

	for name, l := range locs {
		loc, err := models.CreateLocation(db, sc, name, false)
		if err != nil {
			return nil, err
		}
		for i := 0; i < l.NumA; i++ {
			_, err := loc.CreateLocationCard(db, sc, "A")
			if err != nil {
				return nil, err
			}
		}
		letters := []string{"B", "C", "D", "E", "F", "G", "H"}
		for i := 0; i < l.NumOther; i++ {
			_, err := loc.CreateLocationCard(db, sc, letters[i])
			if err != nil {
				return nil, err
			}
		}
		locCards, err := loc.GetCards(db)
		if err != nil {
			return nil, err
		}
		l.Loc = loc
		l.Cards = locCards
	}

	locLinks := []struct {
		Card *models.Card
		Loc  *models.Location
	}{
		{Card: locs["Infirmerie"].Cards[2], Loc: locs["Cabinet"].Loc},
		{Card: locs["Promenade"].Cards[2], Loc: locs["Parc"].Loc},
		{Card: locs["Dortoir"].Cards[5], Loc: locs["Catacombes"].Loc},
		{Card: locs["Cabinet"].Cards[3], Loc: locs["Parc"].Loc},
		{Card: locs["Cabinet"].Cards[4], Loc: locs["Labyrinthe"].Loc},
		{Card: locs["Labyrinthe"].Cards[0], Loc: locs["Parc"].Loc},
		{Card: locs["Parc"].Cards[4], Loc: locs["Porte pentacles"].Loc},
		{Card: locs["Parc"].Cards[3], Loc: locs["Serre"].Loc},
		{Card: locs["Catacombes"].Cards[4], Loc: locs["Porte pentacles"].Loc},
		{Card: locs["Porte pentacles"].Cards[3], Loc: locs["Tombeau"].Loc},
		{Card: locs["Porte pentacles"].Cards[4], Loc: locs["Crypte"].Loc},
		{Card: locs["Porte pentacles"].Cards[5], Loc: locs["Catacombes"].Loc},
	}

	for _, ll := range locLinks {
		_, err := models.CreateLocationLink(db, ll.Card, ll.Loc)
		if err != nil {
			return nil, err
		}
	}

	// TODO state token links
	// TODO skill tests

	return sc, nil
}
