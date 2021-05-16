package commonutil

import (
	"crypto/sha256"
	"encoding/base64"
)

func GenCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}
