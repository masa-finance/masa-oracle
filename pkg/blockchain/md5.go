package blockchain

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 is just syntax sugar around crypto/md5
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
