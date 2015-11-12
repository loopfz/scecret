package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	// TODO icon data
}

func NewIcon(c *gin.Context, in *NewIconIn) (*models.Icon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	// TODO resize
	// TODO upload icon data to cloud storage

	return models.CreateIcon(db, sc, "", "")
}

type ListIconsIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListIcons(c *gin.Context, in *ListIconsIn) ([]*models.Icon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListIcons(db, sc)
}

type GetIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDIcon     int64 `path:"icon, required"`
}

func GetIcon(c *gin.Context, in *GetIconIn) (*models.Icon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadIconFromID(db, sc, in.IDIcon)
}

type UpdateIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDIcon     int64 `path:"icon, required"`
	// TODO icon data
}

func UpdateIcon(c *gin.Context, in *UpdateIconIn) (*models.Icon, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return nil, err
	}

	// TODO delete icon data from cloud storage
	// TODO resize
	// TODO upload new icon data to cloud storage

	err = ico.Update(db, "", "")
	if err != nil {
		return nil, err
	}

	return ico, nil
}

type DeleteIconIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDIcon     int64 `path:"icon, required"`
}

func DeleteIcon(c *gin.Context, in *DeleteIconIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return err
	}

	// TODO delete icon data from cloud storage

	return ico.Delete(db)
}
