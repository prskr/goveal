package routing

import (
	"net/http"
	"strings"
	"time"
)

var epoch = time.Unix(0, 0).Format(time.RFC1123)

var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

func NoCache(h http.Handler, pathsToDisableCache []string) http.Handler {

	pathLookup := make(map[string]bool)

	for idx := range pathsToDisableCache {
		pathLookup[strings.ToLower(pathsToDisableCache[idx])] = true
	}

	fn := func(w http.ResponseWriter, r *http.Request) {

		if _, shouldBeHandled := pathLookup[strings.ToLower(r.URL.Path)]; !shouldBeHandled {
			h.ServeHTTP(w, r)
			return
		}

		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
