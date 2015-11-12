package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewStateTokenLinkIn struct {
	IDScenario      int64 `path:"scenario, required"`
	IDCard          int64 `json:"id_card" binding:"required"`
	IDStateToken    int64 `json:"id_card" binding:"required"`
	UnlocksUnlocked bool  `json:"unlocks_unlocked" binding:"required"`
}

func NewStateTokenLink(c *gin.Context, in *NewStateTokenLinkIn) (*models.StateTokenLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	tk, err := models.LoadStateTokenFromID(db, in.IDStateToken)
	if err != nil {
		return nil, err
	}

	return models.CreateStateTokenLink(db, card, tk, in.UnlocksUnlocked)
}

type ListStateTokenLinksIn struct {
	IDScenario   int64  `path:"scenario, required"`
	IDCard       *int64 `query:"id_card"`
	IDStateToken *int64 `query:"id_state_token"`
}

func ListStateTokenLinks(c *gin.Context, in *ListStateTokenLinksIn) ([]*models.StateTokenLink, error) {

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

	var tk *models.StateToken
	if in.IDStateToken != nil {
		tk, err = models.LoadStateTokenFromID(db, *in.IDStateToken)
		if err != nil {
			return nil, err
		}
	}

	return models.ListStateTokenLinks(db, sc, card, tk)
}

type GetStateTokenLinkIn struct {
	IDScenario       int64 `path:"scenario, required"`
	IDStateTokenLink int64 `path:"statetokenlink, required"`
}

func GetStateTokenLink(c *gin.Context, in *GetStateTokenLinkIn) (*models.StateTokenLink, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadStateTokenLinkFromID(db, sc, in.IDStateTokenLink)
}

type DeleteStateTokenLinkIn struct {
	IDScenario       int64 `path:"scenario, required"`
	IDStateTokenLink int64 `path:"statetokenlink, required"`
}

func DeleteStateTokenLink(c *gin.Context, in *DeleteStateTokenLinkIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	tkl, err := models.LoadStateTokenLinkFromID(db, sc, in.IDStateTokenLink)
	if err != nil {
		return err
	}

	return tkl.Delete(db)
}
