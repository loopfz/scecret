package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewLocationIn struct {
	IDScenario int64  `path:"scenario, required"`
	Name       string `json:"name" binding:"required"`
	Hidden     bool   `json:"hidden" binding:"required"`
}

func NewLocation(c *gin.Context, in *NewLocationIn) (*models.Location, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.CreateLocation(db, sc, in.Name, in.Hidden)
}

type ListLocationsIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListLocations(c *gin.Context, in *ListLocationsIn) ([]*models.Location, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListLocations(db, sc)
}

type GetLocationIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `path:"location, required"`
}

func GetLocation(c *gin.Context, in *GetLocationIn) (*models.Location, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadLocationFromID(db, sc, in.IDLoc)
}

type UpdateLocationIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDLoc      int64  `path:"location, required"`
	Name       string `json:"name" binding:"required"`
	Hidden     bool   `json:"hidden" binding:"required"`
	Notes      string `json:"notes"`
}

func UpdateLocation(c *gin.Context, in *UpdateLocationIn) (*models.Location, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	err = loc.Update(db, in.Name, in.Hidden, in.Notes)
	if err != nil {
		return nil, err
	}

	return loc, nil
}

type DeleteLocationIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `path:"location, required"`
}

func DeleteLocation(c *gin.Context, in *DeleteLocationIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return err
	}

	return loc.Delete(db)
}

type NewLocationCardIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDLoc      int64  `path:"location, required"`
	Letter     string `json:"letter" binding:"required"`
}

func NewLocationCard(c *gin.Context, in *NewLocationCardIn) (*models.LocationCard, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	return loc.CreateLocationCard(db, sc, in.Letter)
}

type ListLocationCardsIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `path:"location, required"`
}

func ListLocationCards(c *gin.Context, in *ListLocationCardsIn) ([]*models.LocationCard, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	return loc.ListLocationCards(db)
}

type GetLocationCardIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `path:"location, required"`
	IDLocCard  int64 `path:"location_card, required"`
}

func GetLocationCard(c *gin.Context, in *GetLocationCardIn) (*models.LocationCard, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	return loc.LoadLocationCardFromID(db, in.IDLocCard)
}

type UpdateLocationCardIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDLoc      int64  `path:"location, required"`
	IDLocCard  int64  `path:"location_card, required"`
	Letter     string `json:"letter" binding:"required"`
}

func UpdateLocationCard(c *gin.Context, in *UpdateLocationCardIn) (*models.LocationCard, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return nil, err
	}

	lc, err := loc.LoadLocationCardFromID(db, in.IDLocCard)
	if err != nil {
		return nil, err
	}

	err = lc.Update(db, in.Letter)
	if err != nil {
		return nil, err
	}

	return lc, nil
}

type DeleteLocationCardIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDLoc      int64 `path:"location, required"`
	IDLocCard  int64 `path:"location_card, required"`
}

func DeleteLocationCard(c *gin.Context, in *DeleteLocationCardIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	loc, err := models.LoadLocationFromID(db, sc, in.IDLoc)
	if err != nil {
		return err
	}

	lc, err := loc.LoadLocationCardFromID(db, in.IDLocCard)
	if err != nil {
		return err
	}

	return lc.Delete(db)
}
