package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewElementIn struct {
	IDScenario  int64  `path:"scenario, required"`
	Number      int    `json:"number" binding:"required"`
	Description string `json:"description"`
}

func NewElement(c *gin.Context, in *NewElementIn) (*models.Element, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.CreateElement(db, sc, in.Number, in.Description)
}

type ListElementsIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListElements(c *gin.Context, in *ListElementsIn) ([]*models.Element, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListElements(db, sc)
}

type GetElementIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDElem     int64 `path:"element, required"`
}

func GetElement(c *gin.Context, in *GetElementIn) (*models.Element, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadElementFromID(db, sc, in.IDElem)
}

type UpdateElementIn struct {
	IDScenario  int64  `path:"scenario, required"`
	IDElem      int64  `path:"element, required"`
	Number      int    `json:"number" binding:"required"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
}

func UpdateElement(c *gin.Context, in *UpdateElementIn) (*models.Element, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	elem, err := models.LoadElementFromID(db, sc, in.IDElem)
	if err != nil {
		return nil, err
	}

	err = elem.Update(db, in.Number, in.Description, in.Notes)
	if err != nil {
		return nil, err
	}

	return elem, nil
}

type DeleteElementIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDElem     int64 `path:"element, required"`
}

func DeleteElement(c *gin.Context, in *DeleteElementIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	elem, err := models.LoadElementFromID(db, sc, in.IDElem)
	if err != nil {
		return err
	}

	return elem.Delete(db)
}
