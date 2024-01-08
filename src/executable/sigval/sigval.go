package sigval

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"binopoy/mexer/constants"
)

func ValidateSignatures(mexName string, pubkey []byte) (success bool) {
	// Validate signature
	unMexedFilePath := constants.MexerTempDir + "/unmexed/" + mexName + "/"
	validSignatureCount := 0

	for _, possibleName := range constants.PossibleExamPackageNames {
		encryptedExam, err := filepath.Abs(unMexedFilePath + possibleName)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(encryptedExam); err != nil {
			continue
		}
		signature := encryptedExam + ".sig"
		if _, err := os.Stat(signature); err != nil {
			continue
		}

		success, err = validateSignature(encryptedExam, signature, pubkey)
		if err != nil {
			log.Fatal(err)
		}

		if success {
			validSignatureCount++
		}
	}

	if validSignatureCount > 0 && validSignatureCount < 3 {
		return true
	}
	return false
}

func validateSignature(
	encryptedExamFilePath string,
	signatureFilePath string,
	pubkey []byte,
) (success bool, err error) {
	// read both files
	examContent, err := os.ReadFile(encryptedExamFilePath)
	if err != nil {
		return false, err
	}

	signatureContent, err := os.ReadFile(signatureFilePath)
	if err != nil {
		return false, err
	}

	publicKey, err := getPubKey(pubkey)

	success, err = verifyWithSHA256AndRSA(
		examContent,
		publicKey,
		string(signatureContent),
	)
	if err != nil {
		fmt.Println(err)
		// Check error type
		if err == rsa.ErrVerification {
			return false, nil
		}
		return false, err
	}

	return success, nil
}

func getPubKey(pubkey []byte) (*rsa.PublicKey, error) {
	// Read the contents of the file

	// Decode the PEM-encoded public key
	block, _ := pem.Decode(pubkey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing public key: %v", err)
	}

	// Assert that the public key is an RSA public key
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an RSA public key")
	}

	return rsaPubKey, nil
}

func verifyWithSHA256AndRSA(
	signedDataBuffer []byte,
	publicKey *rsa.PublicKey,
	signature string,
) (bool, error) {
	hashed := sha256.Sum256(signedDataBuffer)

	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("error decoding signature: %v", err)
	}

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sig)
	if err != nil {
		return false, fmt.Errorf("verification failed: %v", err)
	}

	return true, nil
}
