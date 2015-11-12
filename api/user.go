package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/loopfz/scecret/auth"
	"github.com/loopfz/scecret/models"
)

type RegisterUserIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context, in *RegisterUserIn) (*models.User, error) {
	return models.CreateUser(db, in.Email, in.Password)
}

type AuthIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Auth(c *gin.Context, in *AuthIn) (string, error) {

	u, err := models.LoadUserFromEmail(db, in.Email)
	if err != nil {
		return "", err
	}

	if !u.PasswordEquals(in.Password) {
		return "", errors.New("Bad password")
	}

	tk, err := auth.CreateToken(u)
	if err != nil {
		return "", err
	}

	return tk, nil
}

func GetMe(c *gin.Context) (*models.User, error) {

	return auth.RetrieveTokenUser(db, c)
}
