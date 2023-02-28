package api

import (
	"encoding/hex"
	"errors"
	"hash/fnv"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/Masterminds/sprig/v3"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"code.icb4dc0.de/prskr/goveal/config"
	"code.icb4dc0.de/prskr/goveal/fs"
	"code.icb4dc0.de/prskr/goveal/rendering"
	"code.icb4dc0.de/prskr/goveal/web"
)

var indexTmpl *template.Template

func init() {
	if t, err := template.
		New("index").
		Funcs(sprig.FuncMap()).
		Funcs(map[string]any{
			"fileId": func(fileName string) string {
				h := fnv.New32a()
				return hex.EncodeToString(h.Sum([]byte(path.Base(fileName))))
			},
		}).
		ParseFS(web.WebFS, "*.gohtml"); err != nil {
		panic(err)
	} else {
		indexTmpl = t
	}
}

type Views struct {
	logger     *log.Logger
	cfg        *config.Components
	wdfs       fs.FS
	mdFilepath string
}

func RegisterViews(router *httprouter.Router, logger *log.Logger, wdfs fs.FS, mdFilepath string, cfg *config.Components) {
	p := &Views{
		logger:     logger,
		cfg:        cfg,
		wdfs:       wdfs,
		mdFilepath: mdFilepath,
	}
	router.GET("/", p.IndexPage)
	router.GET("/index.html", p.IndexPage)
	router.GET("/index.htm", p.IndexPage)
	router.GET("/slides", p.RenderedMarkdown)
}

func (p *Views) IndexPage(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "text/html")
	if err := indexTmpl.ExecuteTemplate(writer, "index.gohtml", map[string]any{
		"Reveal":    p.cfg.Reveal,
		"Rendering": p.cfg.Rendering,
	}); err != nil {
		p.logger.Errorf("Failed to render template: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
	}
}

func (p *Views) RenderedMarkdown(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	f, err := p.wdfs.Open(p.mdFilepath)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	defer func() {
		err = errors.Join(err, f.Close())
	}()

	data, err := io.ReadAll(f)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Set("Content-Type", "text/html")
	var rendered []byte
	if rendered, err = rendering.ToHTML(string(data), p.cfg.Rendering); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	} else if _, err = writer.Write(rendered); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
}
