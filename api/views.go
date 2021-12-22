package api

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/multierr"

	"github.com/baez90/goveal/config"
	"github.com/baez90/goveal/fs"
	"github.com/baez90/goveal/rendering"
)

type Views struct {
	cfg        *config.Components
	wdfs       fs.FS
	mdFilepath string
}

func RegisterViews(app *fiber.App, wdfs fs.FS, mdFilepath string, cfg *config.Components) {
	p := &Views{cfg: cfg, wdfs: wdfs, mdFilepath: mdFilepath}
	app.Get("/", p.IndexPage)
	app.Get("/index.html", p.IndexPage)
	app.Get("/index.htm", p.IndexPage)
	app.Get("/slides", p.RenderedMarkdown)
}

func (p *Views) IndexPage(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"Reveal":    p.cfg.Reveal,
		"Rendering": p.cfg.Rendering,
	})
}

func (p *Views) RenderedMarkdown(ctx *fiber.Ctx) (err error) {
	f, err := p.wdfs.Open(p.mdFilepath)
	if err != nil {
		return err
	}
	defer multierr.AppendInvoke(&err, multierr.Close(f))
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var rendered []byte
	if rendered, err = rendering.ToHTML(string(data), p.cfg.Rendering); err != nil {
		return err
	} else if _, err = ctx.Write(rendered); err != nil {
		return err
	}

	ctx.Append(fiber.HeaderContentType, fiber.MIMETextHTML)

	return err
}
