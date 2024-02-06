package crypto

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"binopoy/mexer/constants"
)

var possibleExamPackageNames = [2]string{"exam.xml.bin", "abitti-exam.xml.bin"}

func DecryptFile(filePath string, key []byte, iv []byte) (success bool, examName string) {
	// Base file Name
	baseName := filepath.Base(filePath)
	unMexedFilePath := constants.MexerTempDir + "/unmexed/" + baseName + "/"

	examPackageName := ""
	for _, possibleName := range possibleExamPackageNames {
		examFilePath, err := filepath.Abs(unMexedFilePath + possibleName)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := os.Stat(examFilePath); err == nil {
			examPackageName = possibleName
			break
		}
	}
	if examPackageName == "" {
		return false, ""
	}

	examFilePath, err := filepath.Abs(unMexedFilePath + examPackageName)
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := os.ReadFile(examFilePath)
	if err != nil {
		log.Fatal(err)
	}

	decryptedResult, err := decryptAESCTR(key, iv, fileContent)
	if err != nil {
		log.Fatal("Error decrypting:", err)
	}

	// Check that decrypted file is valid xml
	if !strings.Contains(string(decryptedResult), "<?xml") {
		return false, ""
	}

	// Read exam name from decrypted file     <e:exam-title>Uusi koe 2</e:exam-title> field
	examName = strings.Split(string(decryptedResult), "<e:exam-title>")[1]
	examName = strings.Split(examName, "</e:exam-title>")[0]

	// Sanitate exam name to a one line string
	examName = strings.ReplaceAll(examName, "\n", " ")
	examName = strings.ReplaceAll(examName, "<br/>", " ")
	examName = strings.TrimSpace(examName)

	return true, examName
}
