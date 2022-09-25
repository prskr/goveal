package rendering

import (
	"embed"
	"hash/fnv"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/valyala/bytebufferpool"

	"code.icb4dc0.de/prskr/goveal/rendering/emoji"
)

var (
	//go:embed templates/*.gohtml
	templatesFS embed.FS
	templates   *template.Template
)

func init() {
	templates = template.New("rendering").
		Funcs(sprig.FuncMap()).
		Funcs(template.FuncMap{
			"renderMarkdown": func(md string) template.HTML {
				rr := &RevealRenderer{
					Hash: fnv.New32a(),
				}

				emojis := emoji.NewEmojiParser()
				mdParser := parser.NewWithExtensions(parserExtensions)
				mdParser.Opts.ParserHook = emojis.EmojiParser
				renderer := mdhtml.NewRenderer(mdhtml.RendererOptions{
					Flags:          mdhtml.CommonFlags | mdhtml.HrefTargetBlank,
					RenderNodeHook: rr.RenderHook,
				})

				renderedHTML := markdown.ToHTML([]byte(md), mdParser, renderer)
				//nolint:gosec // template should be esacped
				return template.HTML(renderedHTML)
			},
		})
	var err error
	if templates, err = templates.ParseFS(templatesFS, "templates/*.gohtml"); err != nil {
		panic(err)
	}
}

func renderTemplate(templateName string, data any) (output []byte, err error) {
	buffer := bytebufferpool.Get()
	defer func() {
		buffer.Reset()
		bytebufferpool.Put(buffer)
	}()

	err = templates.ExecuteTemplate(buffer, templateName, data)
	return buffer.Bytes(), err
}
