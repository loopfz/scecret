package securerandom

import "crypto/rand"

func RandomBytes(size uint) ([]byte, error) {
	ret := make([]byte, size)
	_, err := rand.Read(ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
