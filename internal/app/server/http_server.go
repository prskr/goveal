package server

import (
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/baez90/goveal/assets"
	"github.com/baez90/goveal/internal/app/rendering"
	"github.com/baez90/goveal/internal/app/routing"
	"github.com/baez90/goveal/internal/encoding"
)

const (
	markdownFilePath = "/content.md"
)

type (
	Config struct {
		Host         string
		Port         uint16
		MarkdownPath string
		RevealParams *rendering.RevealParams
	}

	HTTPServer struct {
		listener net.Listener
		handler  http.Handler
	}
)

func (srv HTTPServer) Serve() error {
	return http.Serve(srv.listener, srv.handler)
}

func (srv HTTPServer) ListenAddress() string {
	return srv.listener.Addr().String()
}

func NewHTTPServer(config Config) (srv *HTTPServer, err error) {
	noCacheFiles := append(config.RevealParams.FilesToMonitor, markdownFilePath)
	if err := detectMarkdownFileEnding(config.MarkdownPath, config.RevealParams); err != nil {
		return nil, err
	}

	router := &routing.RegexpRouter{}
	var tmplRenderer rendering.RevealRenderer
	if tmplRenderer, err = rendering.NewRevealRenderer(config.RevealParams); err != nil {
		err = fmt.Errorf("failed to initialize reveal renderer %w", err)
		return
	}

	// language=regexp
	if err = router.AddRule(`^(/(index.html(l)?)?)?$`, tmplRenderer); err != nil {
		return
	}

	var mdFS http.FileSystem
	if mdFS, err = routing.NewMarkdownFS(config.MarkdownPath); err != nil {
		err = fmt.Errorf("failed to initialize markdown file handler %w", err)
		return
	}

	var revealFS, webFS fs.FS
	if revealFS, err = fs.Sub(assets.Assets, "reveal"); err != nil {
		return nil, err
	}

	if webFS, err = fs.Sub(assets.Assets, "web"); err != nil {
		return nil, err
	}

	layeredFS := routing.NewLayeredFileSystem(http.FS(revealFS), http.FS(webFS), http.Dir("."), mdFS)

	// language=regexp
	if err = router.AddRule(`^(?i)/hash/(md5|sha1|sha2)/.*`, NewHashHandler(layeredFS)); err != nil {
		return
	}
	// language=regexp
	if err = router.AddRule("^/.*\\.md$", http.FileServer(mdFS)); err != nil {
		return
	}
	// language=regexp
	if err = router.AddRule("/.+", http.FileServer(layeredFS)); err != nil {
		return
	}

	hostPort := fmt.Sprintf("%s:%d", config.Host, config.Port)

	srv = &HTTPServer{
		handler: routing.NoCache(router, noCacheFiles),
	}

	if srv.listener, err = net.Listen("tcp", hostPort); err != nil {
		return
	}

	return
}

func detectMarkdownFileEnding(filePath string, params *rendering.RevealParams) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if le, err := encoding.Detect(f); err != nil {
		return err
	} else {
		params.LineEnding = le
	}
	return nil
}
