package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/qeaml/templd"
)

//go:embed data/templates
var templates embed.FS

func main() {
	templateFS, err := fs.Sub(templates, "data/templates")
	if err != nil {
		log.Fatalln(err)
	}
	views := templd.NewViews(http.FS(templateFS), ".txt").EmbedStyle()
	app := fiber.New(fiber.Config{
		AppName: "CMPD",
		Views:   views,
	})
	app.Get("/test/:templ", func(c *fiber.Ctx) error {
		return c.Render(c.Params("templ"), templd.Vars{})
	})
	log.Fatalln(app.Listen(":1987"))
}
