package encrypt

import (
	"crypto/md5"
	"encoding/hex"
)

func Encrypt(pass string) string {
	hash := md5.Sum([]byte(pass))
	hashedPass := hex.EncodeToString(hash[:])

	return hashedPass
}
