package api

import "github.com/gofiber/fiber/v2"

func NoCache(app *fiber.App) {
	app.Use(NoCacheHandler)
}

func NoCacheHandler(ctx *fiber.Ctx) error {
	ctx.Response().Header.Set("Cache-Control", "no-cache")
	return ctx.Next()
}
