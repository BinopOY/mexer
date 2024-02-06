package main

import (
	_ "embed"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html/v2"

	"binopoy/mexerui/handlers"
)

var executableTag = "v0.1.3-beta"

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views:     engine,
		BodyLimit: 4 * 1024 * 1024 * 1024,
	})

	downloadExecutable()

	app.Get("/", handlers.Index)

	app.Use(cors.New())

	app.Post("/validate", handlers.Validate)

	app.Static("/", "./static")

	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": os.Getenv("BASIC_AUTH_PASSWORD"),
		},
	}))
	app.Get("/metrics", monitor.New())

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}

func downloadExecutable() {
	executableFile, err := os.Create("mexer_amd64")
	if err != nil {
		panic(err)
	}
	defer executableFile.Close()
	resp, err := http.Get(
		"https://github.com/BinopOY/mexer/releases/download/" + executableTag + "/mexer_amd64_linux",
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(executableFile, resp.Body)
	if err != nil {
		panic(err)
	}

	// Make the file executable
	err = os.Chmod("mexer_amd64", 0755)
}
