package api

import (
	"net/http"
)

func NoCache(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Cache-Control", "no-cache")
		handler.ServeHTTP(writer, request)
	})
}
