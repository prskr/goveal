package api

import (
	"hash/fnv"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"go.uber.org/multierr"

	"github.com/baez90/goveal/config"
	"github.com/baez90/goveal/fs"
	"github.com/baez90/goveal/rendering"
)

const (
	parserExtensions = parser.NoIntraEmphasis | parser.Tables | parser.FencedCode |
		parser.Autolink | parser.Strikethrough | parser.SpaceHeadings | parser.HeadingIDs |
		parser.BackslashLineBreak | parser.DefinitionLists | parser.MathJax | parser.Titleblock |
		parser.OrderedListStart | parser.Attributes
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

	mdParser := parser.NewWithExtensions(parserExtensions)
	rr := &rendering.RevealRenderer{
		StateMachine: rendering.NewStateMachine("***", "---"),
		Hash:         fnv.New32a(),
	}
	renderer := html.NewRenderer(html.RendererOptions{
		Flags:          html.CommonFlags | html.HrefTargetBlank,
		RenderNodeHook: rr.RenderHook,
	})

	ctx.Append(fiber.HeaderContentType, fiber.MIMETextHTML)
	if _, err = ctx.Write(markdown.ToHTML(data, mdParser, renderer)); err != nil {
		return err
	}

	return err
}
