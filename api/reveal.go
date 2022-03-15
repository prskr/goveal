package api

import (
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/baez90/goveal/assets"
	"github.com/baez90/goveal/fs"
	"github.com/baez90/goveal/web"
)

func FileSystemMiddleware(fallthroughHandler http.Handler, wdfs fs.FS) http.Handler {
	layers := []fs.FS{wdfs}
	layers = append([]fs.FS{web.WebFS, assets.Assets}, layers...)

	layeredFS := fs.Layered{Layers: layers}
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		reqPath := strings.TrimPrefix(req.URL.Path, "/")
		if f, err := layeredFS.Open(reqPath); err != nil {
			fallthroughHandler.ServeHTTP(writer, req)
			return
		} else if readSeeker, ok := f.(io.ReadSeeker); ok {
			http.ServeContent(writer, req, path.Base(reqPath), time.Now(), readSeeker)
			_ = f.Close()
		} else {
			_ = f.Close()
			fallthroughHandler.ServeHTTP(writer, req)
		}
	})
}
