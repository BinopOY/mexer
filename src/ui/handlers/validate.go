package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var tempDir = "./.mexerui-temp"

type failedExam struct {
	FileName string
	Reason   string
}

func Validate(c *fiber.Ctx) error {
	// Get the file from the request:
	file, err := c.FormFile("f")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Get the text from the request:
	text := c.FormValue("t")

	// Make sure the temp directory exists:
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, 0755)
	}
	defer os.RemoveAll(tempDir)

	// Save file to temp directory:
	if err := c.SaveFile(file, tempDir+"/"+file.Filename); err != nil {
		fmt.Println(err)
		c.SendStatus(500)
		return err
	}

	// Save codes to temp directory:
	if err := os.WriteFile(tempDir+"/"+"codes.txt", []byte(text), 0644); err != nil {
		fmt.Println(err)
		c.SendStatus(500)
		return err
	}

	// Run mexer_amd64 executable:
	out, err := exec.Command("./mexer_amd64", tempDir+"/"+file.Filename, tempDir+"/"+"codes.txt").
		Output()
	if err != nil {
		fmt.Println(err)
		c.SendStatus(500)
		return err
	}

	var successful []string
	var failed []failedExam
	var unused []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "success") {
			examName := strings.ReplaceAll(strings.Split(line, " ")[1], "_", " ")
			successful = append(successful, examName)
		} else if strings.Contains(line, "failed") {
			reason := strings.Split(line, " ")[2]
			finalReason := ""

			if c.Query("json", "false") == "true" {
				finalReason = reason
			} else {
				if reason == "no_code_found" {
					finalReason = "Kokeelle ei löytynyt sopivaa koodia"
				} else if reason == "invalid_signature" {
					finalReason = "Kokeen digitaalinen allekirjoitus ei täsmännyt"
				} else if reason == "unpack" {
					finalReason = "Koepaketin purkaminen epäonnistui"
				}
			}

			failed = append(failed, failedExam{
				FileName: strings.Split(line, " ")[1],
				Reason:   finalReason,
			})
		} else if strings.Contains(line, "unused_code") {
			unused = append(unused, strings.Join(strings.Split(line, " ")[1:], " "))
		}
	}

	if c.Query("json", "false") == "true" {
		return c.JSON(fiber.Map{
			"successes":    successful,
			"failures":     failed,
			"unused_codes": unused,
		})
	}

	return c.Render("partial/result.tmpl", fiber.Map{
		"Successes":   successful,
		"Failures":    failed,
		"UnusedCodes": unused,
	})
}
