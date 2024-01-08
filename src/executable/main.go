package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"binopoy/mexer/constants"
	"binopoy/mexer/crypto"
	"binopoy/mexer/fileops"
	"binopoy/mexer/sigval"
)

type failedExam struct {
	fileName string
	reason   string
}

//go:embed embed-data/pubkey.pem
var pubkey []byte

func main() {
	// flag
	flag.Parse()
	zipPath := flag.Arg(0)
	codePath := flag.Arg(1)

	// check args
	if zipPath == "" || codePath == "" {
		log.Fatal("Usage: ./main <zip file path> <code file path>")
	}

	// Clear temp dir
	defer os.RemoveAll(constants.MexerTempDir)

	// read codes
	codes, err := fileops.ReadCodes(codePath)
	if err != nil {
		log.Fatal(err)
	}

	// unzip
	err = fileops.UnzipFiles(zipPath, constants.MexerTempDir+"/unzipped")
	if err != nil {
		log.Fatal(err)
	}

	files := fileops.GetMexFiles()

	if len(files) == 0 {
		// All files are processed
		log.Fatal("No files to process")
	}

	var failedExams []failedExam

	for _, file := range files {
		// Unmex files
		mexName := filepath.Base(file)
		err = fileops.UnzipFiles(file, constants.MexerTempDir+"/unmexed/"+mexName)
		if err != nil {
			failedExams = append(
				failedExams,
				failedExam{fileName: filepath.Base(file), reason: "unpack"},
			)

			os.Remove(file)
			continue
		}
		// Check file signature
		validSignature := sigval.ValidateSignatures(mexName, pubkey)
		if !validSignature {
			failedExams = append(
				failedExams,
				failedExam{fileName: filepath.Base(file), reason: "invalid_signature"},
			)

			os.Remove(file)
			continue
		}
	}

	var successfulExams []string

	// Decrypt
	for _, code := range codes {
		files = fileops.GetMexFiles()
		if len(files) == 0 {
			// All files are processed
			break
		}
		// Remove all whitespaces from code
		noWhiteSpaceCode := strings.ReplaceAll(code, " ", "")

		key, iv, err := crypto.DeriveAES256KeyAndIV(noWhiteSpaceCode)
		if err != nil {
			log.Fatal(err)
		}
		fileDecrypted := false

		for _, file := range files {
			// Decrypt file
			successful, examName := crypto.DecryptFile(file, key, iv)
			if successful {
				// File decrypted successfully
				// Remove file from unzipped dir
				err = os.Remove(file)
				if err != nil {
					log.Fatal(err)
				}
				successfulExams = append(successfulExams, examName)
				// Break for loop
				fileDecrypted = true
				break
			}
		}
		if !fileDecrypted {
			fmt.Println("unused_code:", code)
		}

	}

	remainingFiles := fileops.GetMexFiles()
	for _, file := range remainingFiles {
		// Push to failed exams
		failedExams = append(
			failedExams,
			failedExam{fileName: filepath.Base(file), reason: "no_code_found"},
		)
	}

	if len(successfulExams) != 0 {
		for _, exam := range successfulExams {
			fmt.Println("success:", strings.ReplaceAll(exam, " ", "_"))
		}
	}

	if len(failedExams) != 0 {
		for _, exam := range failedExams {
			fmt.Println("failed:", strings.ReplaceAll(exam.fileName, " ", "_"), exam.reason)
		}
	}
}
