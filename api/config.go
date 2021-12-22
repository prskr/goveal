package api

import (
	"github.com/gofiber/fiber/v2"

	"github.com/baez90/goveal/config"
)

type ConfigAPI struct {
	cfg *config.Components
}

func RegisterConfigAPI(app *fiber.App, cfg *config.Components) {
	cfgApi := &ConfigAPI{cfg: cfg}
	app.Get("/api/v1/config/reveal", cfgApi.RevealConfig)
	app.Get("/api/v1/config/mermaid", cfgApi.MermaidConfig)
}

func (a *ConfigAPI) RevealConfig(ctx *fiber.Ctx) error {
	return ctx.JSON(a.cfg.Reveal)
}

func (a *ConfigAPI) MermaidConfig(ctx *fiber.Ctx) error {
	return ctx.JSON(a.cfg.Mermaid)
}
