package blockchain

import (
	"encoding/base64"
	"hash"

	"github.com/creachadair/otp"
)

func TOTP(f func() hash.Hash, digits int, t int, key string) string {
	cfg := otp.Config{
		Hash:     f,      // default is sha1.New
		Digits:   digits, // default is 6
		TimeStep: otp.TimeWindow(t),
		Key:      key,
		Format: func(hash []byte, nb int) string {
			return base64.StdEncoding.EncodeToString(hash)[:nb]
		},
	}
	return cfg.TOTP()
}
