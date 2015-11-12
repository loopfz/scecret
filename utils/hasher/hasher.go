package hasher

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hash(s string) string {
	h := sha512.Sum512([]byte(s))
	h2 := h[:]
	return hex.EncodeToString(h2)
}
