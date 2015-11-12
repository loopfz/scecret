package securerandom

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomString(size uint) (string, error) {
	ret := make([]byte, size)
	_, err := rand.Read(ret)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ret), nil
}
