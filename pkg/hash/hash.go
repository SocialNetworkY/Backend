package hash

import (
	"crypto/sha1"
	"encoding/hex"
)

type (
	Config struct {
		Salt string
	}

	SHA1Hasher struct {
		salt string
	}
)

func NewSHA1Hasher(config Config) *SHA1Hasher {
	return &SHA1Hasher{
		salt: config.Salt,
	}
}

func (h *SHA1Hasher) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return hex.EncodeToString(hash.Sum([]byte(h.salt)))
}

func (h *SHA1Hasher) Verify(hash string, password string) bool {
	return h.Hash(password) == hash
}
