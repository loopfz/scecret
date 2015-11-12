package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type ListCardsIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListCards(c *gin.Context, in *ListCardsIn) ([]*models.Card, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListCards(db, sc)
}

type GetCardIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDCard     int64 `path:"card, required"`
}

func GetCard(c *gin.Context, in *GetCardIn) (*models.Card, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadCardFromID(db, sc, in.IDCard)
}

type UpdateCardIn struct {
	IDScenario  int64            `path:"scenario, required"`
	IDCard      int64            `path:"card, required"`
	Number      uint             `json:"number" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Front       *models.CardFace `json:"front" binding:"required"`
	Back        *models.CardFace `json:"back" binding:"required"`
}

func UpdateCard(c *gin.Context, in *UpdateCardIn) (*models.Card, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	err = card.Update(db, in.Number, in.Description, in.Front, in.Back)
	if err != nil {
		return nil, err
	}

	return card, nil
}

type NewCardIconIn struct {
	IDScenario     int64  `path:"scenario, required"`
	IDCard         int64  `path:"card, required"`
	IDIcon         int64  `json:"id_icon" binding:"required"`
	FrontBack      bool   `json:"front_back" binding:"required"`
	X              uint   `json:"x" binding:"required"`
	Y              uint   `json:"y" binding:"required"`
	SizeX          uint   `json:"size_x" binding:"required"`
	SizeY          uint   `json:"size_y" binding:"required"`
	Annotation     string `json:"annotation" binding:"required"`
	AnnotationType int    `json:"annotation_type" binding:"required"`
}

func NewCardIcon(c *gin.Context, in *NewCardIconIn) (*models.CardIcon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return nil, err
	}

	return card.CreateCardIcon(db, ico, in.FrontBack, in.X, in.Y, in.SizeX, in.SizeY,
		in.Annotation, in.AnnotationType, nil, nil)
}

type ListCardIconsIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDCard     int64 `path:"card, required"`
}

func ListCardIcons(c *gin.Context, in *ListCardIconsIn) ([]*models.CardIcon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	return card.ListCardIcons(db, nil, nil)
}

type GetCardIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDCard     int64 `path:"card, required"`
	IDCardIcon int64 `path:"icon, required"`
}

func GetCardIcon(c *gin.Context, in *GetCardIconIn) (*models.CardIcon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	return card.LoadCardIconFromID(db, in.IDCardIcon)
}

type UpdateCardIconIn struct {
	IDScenario     int64  `path:"scenario, required"`
	IDCard         int64  `path:"card, required"`
	IDCardIcon     int64  `path:"icon, required"`
	IDIcon         int64  `json:"id_icon" binding:"required"`
	FrontBack      bool   `json:"front_back" binding:"required"`
	X              uint   `json:"x" binding:"required"`
	Y              uint   `json:"y" binding:"required"`
	SizeX          uint   `json:"size_x" binding:"required"`
	SizeY          uint   `json:"size_y" binding:"required"`
	Annotation     string `json:"annotation" binding:"required"`
	AnnotationType int    `json:"annotation_type" binding:"required"`
}

func UpdateCardIcon(c *gin.Context, in *UpdateCardIconIn) (*models.CardIcon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return nil, err
	}

	ci, err := card.LoadCardIconFromID(db, in.IDCardIcon)
	if err != nil {
		return nil, err
	}

	err = ci.Update(db, ico, in.FrontBack, in.X, in.Y, in.SizeX, in.SizeY,
		in.Annotation, in.AnnotationType)
	if err != nil {
		return nil, err
	}

	return ci, nil
}

type DeleteCardIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDCard     int64 `path:"card, required"`
	IDCardIcon int64 `path:"icon, required"`
}

func DeleteCardIcon(c *gin.Context, in *DeleteCardIconIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return err
	}

	ci, err := card.LoadCardIconFromID(db, in.IDCardIcon)
	if err != nil {
		return err
	}

	return ci.Delete(db)
}
