package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type CreateSkillTestIn struct {
	IDScenario     int64 `path:"scenario, required"`
	IDCard         int64 `json:"id_card" binding:"required"`
	IDStat         int64 `json:"id_stat" binding:"required"`
	NormalShields  uint  `json:"normal_shields"`
	SkullShields   uint  `json:"skull_shields"`
	HeartShields   uint  `json:"heart_shields"`
	UTShields      uint  `json:"ut_shields"`
	SpecialShields uint  `json:"special_shields"`
}

func CreateSkillTest(c *gin.Context, in *CreateSkillTestIn) (*models.SkillTest, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, in.IDCard)
	if err != nil {
		return nil, err
	}

	stat, err := models.LoadStatFromID(db, sc, in.IDStat)
	if err != nil {
		return nil, err
	}

	return models.CreateSkillTest(db, card, stat, in.NormalShields, in.SkullShields,
		in.HeartShields, in.UTShields, in.SpecialShields)
}

type ListSkillTestsIn struct {
	IDScenario int64  `path:"scenario, required"`
	IDCard     *int64 `query:"id_card"`
	IDStat     *int64 `query:"id_stat"`
}

func ListSkillTests(c *gin.Context, in *ListSkillTestsIn) ([]*models.SkillTest, error) {

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

	var stat *models.Stat
	if in.IDStat != nil {
		stat, err = models.LoadStatFromID(db, sc, *in.IDStat)
		if err != nil {
			return nil, err
		}
	}

	return models.ListSkillTests(db, sc, card, stat)
}

type GetSkillTestIn struct {
	IDScenario  int64 `path:"scenario, required"`
	IDSkillTest int64 `path:"skilltest, required"`
}

func GetSkillTest(c *gin.Context, in *GetSkillTestIn) (*models.SkillTest, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadSkillTestFromID(db, sc, in.IDSkillTest)
}

type UpdateSkillTestIn struct {
	IDScenario     int64 `path:"scenario, required"`
	IDSkillTest    int64 `path:"skilltest, required"`
	IDStat         int64 `json:"id_stat" binding:"required"`
	NormalShields  uint  `json:"normal_shields"`
	SkullShields   uint  `json:"skull_shields"`
	HeartShields   uint  `json:"heart_shields"`
	UTShields      uint  `json:"ut_shields"`
	SpecialShields uint  `json:"special_shields"`
}

func UpdateSkillTest(c *gin.Context, in *UpdateSkillTestIn) (*models.SkillTest, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	st, err := models.LoadSkillTestFromID(db, sc, in.IDSkillTest)
	if err != nil {
		return nil, err
	}

	card, err := models.LoadCardFromID(db, sc, st.IDCard)
	if err != nil {
		return nil, err
	}

	stat, err := models.LoadStatFromID(db, sc, in.IDStat)
	if err != nil {
		return nil, err
	}

	err = st.Update(db, card, stat, in.NormalShields, in.SkullShields, in.HeartShields,
		in.UTShields, in.SpecialShields)

	if err != nil {
		return nil, err
	}

	return st, nil
}

type DeleteSkillTestIn struct {
	IDScenario  int64 `path:"scenario, required"`
	IDSkillTest int64 `path:"skilltest, required"`
}

func DeleteSkillTest(c *gin.Context, in *DeleteSkillTestIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	st, err := models.LoadSkillTestFromID(db, sc, in.IDSkillTest)
	if err != nil {
		return err
	}

	return st.Delete(db)
}
