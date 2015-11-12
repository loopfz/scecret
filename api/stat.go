package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewStatIn struct {
	IDScenario  int64  `path:"scenario, required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	IDIcon      int64  `json:"id_icon" binding:"required"`
}

func NewStat(c *gin.Context, in *NewStatIn) (*models.Stat, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return nil, err
	}

	return models.CreateStat(db, sc, ico, in.Name, in.Description)
}

type ListStatsIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListStats(c *gin.Context, in *ListStatsIn) ([]*models.Stat, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListStats(db, sc)
}

type GetStatIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDStat     int64 `path:"stat, required"`
}

func GetStat(c *gin.Context, in *GetStatIn) (*models.Stat, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadStatFromID(db, sc, in.IDStat)
}

type UpdateStatIn struct {
	IDScenario  int64  `path:"scenario, required"`
	IDStat      int64  `path:"stat, required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	IDIcon      int64  `json:"id_icon" binding:"required"`
}

func UpdateStat(c *gin.Context, in *UpdateStatIn) (*models.Stat, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	ico, err := models.LoadIconFromID(db, sc, in.IDIcon)
	if err != nil {
		return nil, err
	}

	st, err := models.LoadStatFromID(db, sc, in.IDStat)
	if err != nil {
		return nil, err
	}

	err = st.Update(db, ico, in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	return st, nil
}

type DeleteStatIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDStat     int64 `path:"stat, required"`
}

func DeleteStat(c *gin.Context, in *DeleteStatIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	st, err := models.LoadStatFromID(db, sc, in.IDStat)
	if err != nil {
		return err
	}

	return st.Delete(db)
}
