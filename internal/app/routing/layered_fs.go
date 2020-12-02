package routing

import (
	"errors"
	"net/http"
	"path/filepath"
)

type LayeredHandler struct {
	layers []http.FileSystem
}

func (l *LayeredHandler) AddHandlers(layers ...http.FileSystem) {
	l.layers = append(l.layers, layers...)
}

func (l LayeredHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_, fileName := filepath.Split(request.URL.Path)
	f, err := l.selectLayer(request.URL.Path)
	if err != nil {
		writer.WriteHeader(404)
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		writer.WriteHeader(500)
		return
	}
	http.ServeContent(writer, request, fileName, stat.ModTime(), f)
}

func (l LayeredHandler) selectLayer(requestedFilePath string) (f http.File, err error) {
	for idx := range l.layers {
		layer := l.layers[idx]
		if f, err = layer.Open(requestedFilePath); err == nil {
			return
		}
	}
	return nil, errors.New("not found")
}
