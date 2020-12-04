package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/baez90/go-reveal-slides/internal/app/rendering"
	"github.com/baez90/go-reveal-slides/internal/app/routing"
	"github.com/markbates/pkger"
)

type Config struct {
	Host         string
	Port         uint16
	MarkdownPath string
	RevealParams *rendering.RevealParams
}

type HTTPServer struct {
	listener net.Listener
	handler  http.Handler
}

func (srv HTTPServer) Serve() error {
	return http.Serve(srv.listener, srv.handler)
}

func (srv HTTPServer) ListenAddress() string {
	return srv.listener.Addr().String()
}

func NewHTTPServer(config Config) (srv *HTTPServer, err error) {
	router := &routing.RegexpRouter{}
	var tmplRenderer rendering.RevealRenderer
	if tmplRenderer, err = rendering.NewRevealRenderer(config.RevealParams); err != nil {
		err = fmt.Errorf("failed to initialize reveal renderer %w", err)
		return
	}

	if err = router.AddRule(`^(\/(index.html(l)?)?)?$`, tmplRenderer); err != nil {
		return
	}

	var mdFS http.FileSystem
	if mdFS, err = routing.NewMarkdownFS(config.MarkdownPath); err != nil {
		err = fmt.Errorf("failed to initialize markdown file handler %w", err)
		return
	}
	fs := routing.NewLayeredFileSystem(pkger.Dir("/assets/reveal"), pkger.Dir("/assets/web"), http.Dir("."), mdFS)

	//language=regexp
	if err = router.AddRule(`^(?i)/hash/(md5|sha1|sha2)/.*`, routing.NoCache(NewHashHandler(fs))); err != nil {
		return
	}
	if err = router.AddRule("^/.*\\.md$", routing.NoCache(http.FileServer(mdFS))); err != nil {
		return
	}
	if err = router.AddRule("/.+", routing.NoCache(http.FileServer(fs))); err != nil {
		return
	}

	hostPort := fmt.Sprintf("%s:%d", config.Host, config.Port)

	srv = &HTTPServer{
		handler: router,
	}

	if srv.listener, err = net.Listen("tcp", hostPort); err != nil {
		return
	}

	return
}
