package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewElementLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDCard     int64 `json:"id_card" binding:"required"`
	IDElem     int64 `json:"id_element" binding:"required"`
	GivesUses  bool  `json:"gives_uses" binding:"required"`
}

func NewElementLink(c *gin.Context, in *NewElementLinkIn) (*models.ElementLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	elem, err := models.LoadElementFromID(db, sc, in.IDElem)
	if err != nil {
		return nil, err
	}

	return models.CreateElementLink(db, card, elem, in.GivesUses)
}

type ListElementLinksIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDCard     *int64 `query:"id_card"`
	IDElem     *int64 `query:"id_element"`
}

func ListElementLinks(c *gin.Context, in *ListElementLinksIn) ([]*models.ElementLink, error) {

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

	var elem *models.Element
	if in.IDElem != nil {
		elem, err = models.LoadElementFromID(db, sc, *in.IDElem)
		if err != nil {
			return nil, err
		}
	}

	return models.ListElementLinks(db, sc, card, elem)
}

type GetElementLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDElem     int64 `path:"elementlink, required"`
}

func GetElementLink(c *gin.Context, in *GetElementLinkIn) (*models.ElementLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadElementLinkFromID(db, sc, in.IDElem)
}

type DeleteElementLinkIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDElem     int64 `path:"elementlink, required"`
}

func DeleteElementLink(c *gin.Context, in *DeleteElementLinkIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	elem, err := models.LoadElementLinkFromID(db, sc, in.IDElem)
	if err != nil {
		return err
	}

	return elem.Delete(db)
}
