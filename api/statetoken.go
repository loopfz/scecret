package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type ListStateTokensIn struct {
	IDScenario int64 `path:"scenario, required"`
}

func ListStateTokens(c *gin.Context, in *ListStateTokensIn) ([]*models.StateToken, error) {

	_, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.ListStateTokens(db)
}

type GetStateTokenIn struct {
	IDScenario int64 `path:"scenario, required"`
	IDTk       int64 `path:"statetoken, required"`
}

func GetStateToken(c *gin.Context, in *GetStateTokenIn) (*models.StateToken, error) {

	_, err := auth.RetrieveTokenScenario(db, c, in.IDScenario)
	if err != nil {
		return nil, err
	}

	return models.LoadStateTokenFromID(db, in.IDTk)
}
