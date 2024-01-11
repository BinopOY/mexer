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
	examPackageArg := flag.Arg(0)
	codeArg := flag.Arg(1)

	// check args
	if examPackageArg == "" || codeArg == "" {
		log.Fatal("Usage: ./main <zip file path> <code file path>")
	}

	// Clear temp dir
	defer os.RemoveAll(constants.MexerTempDir)

	var codes []string
	var err error
	// read codes
	if strings.HasSuffix(codeArg, ".txt") {
		readCodes, err := fileops.ReadCodes(codeArg)
		if err != nil {
			log.Fatal(err)
		}
		codes = readCodes
	} else {
		codes = append(codes, codeArg)
	}

	if len(codes) == 0 {
		// No codes found
		log.Fatal("No codes found in the input")
	}

	// If zip path ends with .mex, then it is a single file
	if strings.HasSuffix(examPackageArg, ".mex") {
		// Copy file to temp dir
		err = fileops.CopyFile(
			examPackageArg,
			fmt.Sprintf("%s/unzipped/%s", constants.MexerTempDir, filepath.Base(examPackageArg)),
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// unzip
		err = fileops.UnzipFiles(examPackageArg, fmt.Sprintf("%s/unzipped", constants.MexerTempDir))
		if err != nil {
			log.Fatal(err)
		}
	}
	files := fileops.GetMexFiles()

	if len(files) == 0 {
		// All files are processed
		log.Fatal("No files to process")
	}

	// Result arrays
	var failedExams []failedExam
	var successfulExams []string

	for _, file := range files {
		// Unmex files = unzip, mex is just a renamed zip file
		mexName := filepath.Base(file)
		err = fileops.UnzipFiles(
			file,
			fmt.Sprintf("%s/unmexed/%s", constants.MexerTempDir, mexName),
		)
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

	files = fileops.GetMexFiles()
	for _, file := range files {
		if len(files) == 0 {
			// All files are processed
			break
		}
		for _, code := range codes {
			// Remove all whitespaces from code
			noWhitespaceCode := strings.ReplaceAll(code, " ", "")
			key, iv, err := crypto.DeriveAES256KeyAndIV(noWhitespaceCode)
			if err != nil {
				log.Fatal(err)
			}
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
				break
			}
		}
	}
	// Decrypt
	/*
		for _, code := range codes {
			files = fileops.GetMexFiles()
			if len(files) == 0 {
				// All files are processed
				break
			}
			// Remove all whitespaces from code
			noWhitespaceCode := strings.ReplaceAll(code, " ", "")

			key, iv, err := crypto.DeriveAES256KeyAndIV(noWhitespaceCode)
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
	*/

	remainingFiles := fileops.GetMexFiles()
	for _, file := range remainingFiles {
		// Push to failed exams
		failedExams = append(
			failedExams,
			failedExam{fileName: filepath.Base(file), reason: "no_code_found"},
		)
	}

	for _, exam := range successfulExams {
		fmt.Println("success:", strings.ReplaceAll(exam, " ", "_"))
	}

	for _, exam := range failedExams {
		fmt.Println("failed:", strings.ReplaceAll(exam.fileName, " ", "_"), exam.reason)
	}
}
