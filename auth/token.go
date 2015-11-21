package auth

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-gorp/gorp"
	"github.com/juju/errors"
	"github.com/loopfz/scecret/models"
	"github.com/loopfz/scecret/utils/hasher"
	"github.com/loopfz/scecret/utils/securerandom"
)

const (
	TOKEN_LEN    = 64
	TOKEN_HEADER = "X-Auth-Token"
)

var (
	tokenStore = make(map[string]string)
	lock       sync.RWMutex
)

func CreateToken(user *models.User) (string, error) {

	tk, err := securerandom.RandomString(TOKEN_LEN)
	if err != nil {
		return "", err
	}

	lock.Lock()
	tokenStore[hasher.Hash(tk)] = user.Email
	lock.Unlock()

	return tk, nil
}

func RetrieveTokenUser(db *gorp.DbMap, c *gin.Context) (*models.User, error) {

	tk := c.Request.Header.Get(TOKEN_HEADER)

	lock.RLock()
	email, ok := tokenStore[hasher.Hash(tk)]
	lock.RUnlock()
	if !ok {
		return nil, errors.NewUnauthorized(nil, "Bad token")
	}

	u, err := models.LoadUserFromEmail(db, email)
	if err != nil {
		return nil, errors.Wrap(err, errors.New("Error retrieving user information"))
	}

	return u, nil
}

func RetrieveTokenScenario(db *gorp.DbMap, c *gin.Context, IDScenario int64) (*models.Scenario, error) {

	u, err := RetrieveTokenUser(db, c)
	if err != nil {
		return nil, err
	}

	sc, err := models.LoadScenarioFromID(db, u, IDScenario)
	if err != nil {
		return nil, err
	}

	if sc.IDAuthor != u.ID {
		return nil, errors.NewNotFound(nil, "No such scenario")
	}

	return sc, nil
}
