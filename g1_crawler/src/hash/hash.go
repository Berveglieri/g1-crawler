package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func HashUrl(url string) string {
	hash := md5.New()
	hash.Write([]byte(url))
	return hex.EncodeToString(hash.Sum(nil))
}
