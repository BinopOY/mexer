package crypto

import (
	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keySize    = 32
	ivSize     = 16
	iterations = 2000
)

func DeriveAES256KeyAndIV(password string) (key, iv []byte, err error) {
	// Perform PBKDF2 key derivation
	derivedData := pbkdf2.Key(
		[]byte(password),
		[]byte(password),
		iterations,
		keySize+ivSize,
		sha1.New,
	)

	// Extract key and IV
	key = derivedData[:keySize]
	iv = derivedData[keySize:]

	return key, iv, nil
}
