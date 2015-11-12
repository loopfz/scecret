package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type NewScenarioIn struct {
	Name string `json:"name" binding:"required"`
}

func NewScenario(c *gin.Context, in *NewScenarioIn) (*models.Scenario, error) {

	u, err := auth.RetrieveTokenUser(db, c)
	if err != nil {
		return nil, err
	}

	s, err := models.CreateScenario(db, in.Name, u)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func ListScenarios(c *gin.Context) ([]*models.Scenario, error) {

	u, err := auth.RetrieveTokenUser(db, c)
	if err != nil {
		return nil, err
	}
	fmt.Printf("User: %s\n", u.Email)

	return models.ListScenarios(db, u)
}

type GetScenarioIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func GetScenario(c *gin.Context, in *GetScenarioIn) (*models.Scenario, error) {

	return auth.RetrieveTokenScenario(db, c, in.IDScenario)
}

type UpdateScenarioIn struct {
	IDScenario int64  `path:"scenario, required"`
	Name       string `json:"name" binding:"required"`
}

func UpdateScenario(c *gin.Context, in *UpdateScenarioIn) (*models.Scenario, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	err = sc.Update(db, in.Name)
	if err != nil {
		return nil, err
	}

	return sc, nil
}

type DeleteScenarioIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func DeleteScenario(c *gin.Context, in *DeleteScenarioIn) error {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return err
	}

	err = sc.Delete(db)
	if err != nil {
		return err
	}

	return nil
}

type GetGraphIn struct {
	IDScenario int64 `path:"scenario,required"`
}

func GetGraph(c *gin.Context, in *GetGraphIn) (interface{}, error) {

	sc, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.Graph(db, sc)
}
