package fs

import (
	"io/fs"
	"os"
)

type FS = fs.FS
type File = fs.File

var Dir = os.DirFS

type Layered struct {
	Layers []FS
}

func (l Layered) Open(name string) (file fs.File, err error) {
	for idx := range l.Layers {
		if file, err = l.Layers[idx].Open(name); err == nil {
			return file, nil
		}
	}
	return nil, err
}
