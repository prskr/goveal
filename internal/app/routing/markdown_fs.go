package routing

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type markdownFs struct {
	destinationPath string
}

func NewMarkdownFS(path string) (fs http.FileSystem, err error) {
	var info os.FileInfo
	info, err = os.Stat(path)
	if err != nil {
		return
	}

	if info.IsDir() || filepath.Ext(info.Name()) != ".md" {
		err = fmt.Errorf("path %s did not pass sanity checks for markdown files", path)
		return
	}

	return &markdownFs{
		destinationPath: path,
	}, nil
}

func (m markdownFs) Open(_ string) (http.File, error) {
	return os.Open(m.destinationPath)
}
