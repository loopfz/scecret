package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewLocationLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `json:"id_location" binding:"required"`
	IDCard     int64 `json:"id_card" binding:"required"`
}

func NewLocationLink(c *gin.Context, in *NewLocationLinkIn) (*models.LocationLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	return models.CreateLocationLink(db, card, loc)
}

type ListLocationLinksIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDCard     *int64 `query:"id_card"`
	IDLoc      *int64 `query:"id_location"`
}

func ListLocationLinks(c *gin.Context, in *ListLocationLinksIn) ([]*models.LocationLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	var card *models.Card
	if in.IDCard != nil {
		card, err = models.LoadCardFromID(db, sc, *in.IDCard)
		if err != nil {
			return nil, err
		}
	}

	var loc *models.Location
	if in.IDLoc != nil {
		fmt.Printf("ID LOC: %d\n", *in.IDLoc)
		loc, err = models.LoadLocationFromID(db, sc, *in.IDLoc)
		if err != nil {
			return nil, err
		}
	}

	return models.ListLocationLinks(db, sc, card, loc)
}

type GetLocationLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLocLink  int64 `path:"locationlink, required"`
}

func GetLocationLink(c *gin.Context, in *GetLocationLinkIn) (*models.LocationLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadLocationLinkFromID(db, sc, in.IDLocLink)
}

type DeleteLocationLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLocLink  int64 `path:"locationlink, required"`
}

func DeleteLocationLink(c *gin.Context, in *DeleteLocationLinkIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	ll, err := models.LoadLocationLinkFromID(db, sc, in.IDLocLink)
	if err != nil {
		return err
	}

	return ll.Delete(db)
}
