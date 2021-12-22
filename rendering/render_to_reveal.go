package rendering

import (
	"bytes"
	"embed"
	"encoding/hex"
	"hash"
	"html"
	"html/template"
	"io"
	"path"
	"regexp"
	"sync"

	"github.com/Masterminds/sprig/v3"
	"github.com/gomarkdown/markdown/ast"
)

var (
	//go:embed templates/*.gohtml
	templatesFS              embed.FS
	templates                *template.Template
	templateRenderBufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
	htmlElementAttributesRegexp = regexp.MustCompile(`(?P<key>[a-z]+(-[a-z]+)*)="(?P<value>.+)"`)
	notesRegexp                 = regexp.MustCompile(`^(?i)notes?`)
)

const (
	mermaidCodeBlock = "mermaid"
)

func init() {
	var err error
	templates = template.New("rendering").Funcs(sprig.FuncMap())
	if templates, err = templates.ParseFS(templatesFS, "templates/*.gohtml"); err != nil {
		panic(err)
	}
}

type RevealRenderer struct {
	StateMachine *StateMachine
	Hash         hash.Hash
	hasNotes     bool
}

func (r *RevealRenderer) RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch b := node.(type) {
	case *ast.Document:
		if entering {
			_, _ = w.Write([]byte("<section>\n"))
			if next := peekNextRuler(b); next != nil && string(next.Literal) == "***" {
				_, _ = w.Write([]byte("<section>\n"))
				r.StateMachine.CurrentState = StateTypeNested
			}
		} else {
			_, _ = w.Write([]byte("</section>\n"))
		}
		return ast.GoToNext, false
	case *ast.ListItem:
		if entering {
			return r.handleListItem(w, b)
		}
		return ast.GoToNext, false
	case *ast.Text:
		if !entering {
			return ast.GoToNext, false
		}
		if notesRegexp.Match(b.Literal) {
			_, err := w.Write([]byte(`<aside class="notes">`))
			r.hasNotes = true
			return ast.SkipChildren, err == nil
		}
		return ast.GoToNext, false
	case *ast.CodeBlock:
		if entering {
			return r.handleCodeBlock(w, b)
		}
		return ast.GoToNext, false
	case *ast.HorizontalRule:
		next := peekNextRuler(b)
		var input = string(b.Literal)
		if next != nil {
			input += string(next.Literal)
		}
		if r.hasNotes {
			_, _ = w.Write([]byte(`</aside>`))
			r.hasNotes = false
		}
		_, _ = w.Write(r.StateMachine.Accept(input))
		return ast.GoToNext, true
	case *ast.Image:
		if entering {
			return r.handleImage(w, b)
		}
		return ast.GoToNext, false
	default:
		return ast.GoToNext, false
	}
}

func (r *RevealRenderer) handleCodeBlock(w io.Writer, code *ast.CodeBlock) (ast.WalkStatus, bool) {
	code.Info = bytes.ToLower(code.Info)
	switch string(code.Info) {
	case mermaidCodeBlock:
		output, err := renderCodeTemplate("mermaid.gohtml", code)
		if err != nil {
			return ast.GoToNext, false
		}
		_, err = w.Write(output)
		return ast.GoToNext, err == nil
	default:
		output, err := renderCodeTemplate("any-code.gohtml", code)
		if err != nil {
			return ast.GoToNext, false
		}
		_, err = w.Write(output)
		return ast.GoToNext, err == nil
	}
}

func (r *RevealRenderer) handleListItem(w io.Writer, listItem *ast.ListItem) (ast.WalkStatus, bool) {
	for _, child := range listItem.Children {
		if p, ok := child.(*ast.Paragraph); ok {
			if len(p.Children) == 0 {
				return ast.GoToNext, false
			}

			data := map[string]interface{}{
				"Attributes": getAttributesFromChildSpan(p),
			}

			if rendered, err := renderTemplate("listItem.gohtml", data); err != nil {
				return ast.GoToNext, false
			} else if _, err = w.Write(rendered); err != nil {
				return ast.GoToNext, false
			}

			return ast.GoToNext, true
		}
	}
	return ast.GoToNext, false
}

func (r *RevealRenderer) handleImage(w io.Writer, img *ast.Image) (ast.WalkStatus, bool) {
	var title string
	if len(img.Children) >= 1 {
		if txt, ok := img.Children[0].(*ast.Text); ok {
			title = string(txt.Literal)
		}
	}

	data := map[string]interface{}{
		"ID":              hex.EncodeToString(r.Hash.Sum([]byte(path.Base(string(img.Destination))))),
		"Attributes":      getAttributesFromChildSpan(img.GetParent()),
		"ImageSource":     string(img.Destination),
		"AlternativeText": html.EscapeString(title),
	}

	if rendered, err := renderTemplate("image.gohtml", data); err != nil {
		return ast.GoToNext, false
	} else if _, err = w.Write(rendered); err != nil {
		return ast.GoToNext, false
	}

	return ast.SkipChildren, true
}

func getAttributesFromChildSpan(node ast.Node) []template.HTMLAttr {
	if getChildren, ok := node.(interface{ GetChildren() []ast.Node }); ok {
		childs := getChildren.GetChildren()
		if len(childs) == 0 {
			return nil
		}
		for idx := range childs {
			if span, ok := childs[idx].(*ast.HTMLSpan); ok {
				return extractElementAttributes(span)
			}
		}
	}
	return nil
}

func extractElementAttributes(htmlSpan *ast.HTMLSpan) (attrs []template.HTMLAttr) {
	htmlComment := string(htmlSpan.Literal)
	if htmlComment == "" {
		return nil
	}
	matches := htmlElementAttributesRegexp.FindAllStringSubmatch(htmlComment, -1)
	attrs = make([]template.HTMLAttr, 0, len(matches))
	for idx := range matches {
		if len(matches[idx]) != 4 {
			continue
		}

		//nolint:gosec // it's the user's responsibility to not skrew this up here
		attrs = append(attrs, template.HTMLAttr(matches[idx][0]))
	}

	return attrs
}

func renderCodeTemplate(templateName string, codeBlock *ast.CodeBlock) (output []byte, err error) {
	data := map[string]interface{}{
		//nolint:gosec // need to embed the code in original format without escaping
		"Code":        template.HTML(codeBlock.Literal),
		"LineNumbers": lineNumbers(codeBlock.Attribute),
	}

	return renderTemplate(templateName, data)
}

func renderTemplate(templateName string, data interface{}) (output []byte, err error) {
	buffer := templateRenderBufferPool.Get().(*bytes.Buffer)
	defer func() {
		buffer.Reset()
		templateRenderBufferPool.Put(buffer)
	}()

	err = templates.ExecuteTemplate(buffer, templateName, data)
	return buffer.Bytes(), err
}

func lineNumbers(attrs *ast.Attribute) string {
	if attrs == nil || attrs.Attrs == nil {
		return ""
	}
	return string(attrs.Attrs["line-numbers"])
}
