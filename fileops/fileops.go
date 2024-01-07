package fileops

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"binopoy/mexer/constants"
)

func UnzipFiles(zipPath string, destination string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	destination, err = filepath.Abs(destination)
	if err != nil {
		return err
	}

	var fileNames []string

	for _, f := range reader.File {
		err := unzipFile(f, destination)
		if err != nil {
			return err
		}
		fileNames = append(fileNames, f.FileInfo().Name())
	}
	return nil
}

func unzipFile(f *zip.File, destination string) error {
	filePath := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path: %s", filePath)
	}

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// 6. Create a destination file for unzipped content
	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 7. Unzip the content of a file and copy it to the destination file
	zippedFile, err := f.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	if _, err := io.Copy(destinationFile, zippedFile); err != nil {
		return err
	}
	return nil
}

func ReadCodes(codePath string) ([]string, error) {
	// Read code file and parse every line as a code
	codeFile, err := os.Open(codePath)
	if err != nil {
		log.Fatal(err)
	}
	defer codeFile.Close()
	var codes []string
	scanner := bufio.NewScanner(codeFile)
	for scanner.Scan() {
		code := scanner.Text()
		codes = append(codes, code)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return codes, nil
}

func GetMexFiles() []string {
	var files []string
	err := filepath.Walk(
		constants.MexerTempDir+"/unzipped",
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			files = append(files, path)
			return nil
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return files
}
