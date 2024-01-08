package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func decryptAESCTR(encKey, iv, fileToDecrypt []byte) ([]byte, error) {
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	mode := cipher.NewCTR(block, iv)
	decryptedResult := make([]byte, len(fileToDecrypt))
	mode.XORKeyStream(decryptedResult, fileToDecrypt)

	return decryptedResult, nil
}
