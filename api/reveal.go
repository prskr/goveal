package api

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/baez90/goveal/assets"
	"github.com/baez90/goveal/fs"
	"github.com/baez90/goveal/web"
)

func RegisterStaticFileHandling(app *fiber.App, wdfs fs.FS) error {
	layers := []fs.FS{wdfs}
	layers = append([]fs.FS{web.WebFS, assets.Assets}, layers...)

	layeredFS := fs.Layered{Layers: layers}
	fsMiddleware := filesystem.New(filesystem.Config{
		Root: http.FS(layeredFS),
		Next: func(c *fiber.Ctx) bool {
			_, err := layeredFS.Open(strings.TrimLeft(c.Path(), "/"))
			return err != nil
		},
	})

	app.Use(fsMiddleware)

	return nil
}
