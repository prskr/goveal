package server

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"hash"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"go.uber.org/multierr"
)

var (
	pathMatcherRegexp = regexp.MustCompile(`(?i)^/hash/(md5|sha1|sha2)(/.*)`)
	hashes            = map[string]crypto.Hash{
		"md5":    crypto.MD5,
		"sha1":   crypto.SHA1,
		"sha256": crypto.SHA256,
	}
)

type hashResponse struct {
	FilePath string
	Hash     string
}

func NewHashHandler(fs http.FileSystem) http.Handler {
	return &hashHandler{
		fs: fs,
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 4096)
			},
		},
	}
}

type hashHandler struct {
	bufferPool *sync.Pool
	fs         http.FileSystem
}

func (h hashHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	components := pathMatcherRegexp.FindStringSubmatch(request.URL.Path)
	if len(components) != 3 {
		writer.WriteHeader(400)
		return
	}

	filePath := components[2]

	var hashFound bool
	var hashInstance crypto.Hash
	if hashInstance, hashFound = hashes[strings.ToLower(components[1])]; !hashFound {
		writer.WriteHeader(404)
		return
	}
	var f http.File
	var err error
	if f, err = h.fs.Open(filePath); err != nil {
		writer.WriteHeader(404)
		return
	}

	var encodedHash string
	if encodedHash, err = h.buildHash(hashInstance.New(), f); err != nil {
		log.Errorf("Failed to calculate hash %v", err)
		writer.WriteHeader(500)
		return
	}

	resp := hashResponse{
		FilePath: filePath,
		Hash:     encodedHash,
	}

	encoder := json.NewEncoder(writer)
	if err = encoder.Encode(resp); err != nil {
		writer.WriteHeader(500)
		return
	}
}

func (h *hashHandler) buildHash(hasher hash.Hash, src io.ReadCloser) (encodedHash string, err error) {
	defer func() {
		err = multierr.Append(err, src.Close())
	}()

	buffer := h.bufferPool.Get().([]byte)
	for n, err := src.Read(buffer); n > 0 && err == nil; n, err = src.Read(buffer) {
		if _, err := hasher.Write(buffer[:n]); err != nil {
			return "", err
		}
	}
	encodedHash = hex.EncodeToString(hasher.Sum(nil))
	return
}
