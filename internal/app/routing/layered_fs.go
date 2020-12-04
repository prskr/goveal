package routing

import (
	"errors"
	"net/http"
	"sync"
)

var (
	ErrFileNotFound = errors.New("file not found in any layer")
)

func NewLayeredFileSystem(layers ...http.FileSystem) http.FileSystem {
	return &layeredFileSystem{
		resolveCache: make(map[string]http.FileSystem),
		layers:       layers,
	}
}

type layeredFileSystem struct {
	layers       []http.FileSystem
	resolveCache map[string]http.FileSystem
	lock         sync.Mutex
}

func (l *layeredFileSystem) Open(name string) (f http.File, err error) {
	if cachedLayer, isCached := l.resolveCache[name]; isCached {
		if cachedLayer != nil {
			return cachedLayer.Open(name)
		} else {
			return nil, ErrFileNotFound
		}
	}

	for idx := range l.layers {
		layer := l.layers[idx]
		if f, err = layer.Open(name); err == nil {
			l.lock.Lock()
			l.resolveCache[name] = layer
			l.lock.Unlock()
			return
		}
	}

	l.lock.Lock()
	l.resolveCache[name] = nil
	l.lock.Unlock()
	return nil, ErrFileNotFound
}
