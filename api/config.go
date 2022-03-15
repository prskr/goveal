package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/baez90/goveal/config"
)

type ConfigAPI struct {
	cfg *config.Components
}

func RegisterConfigAPI(router *httprouter.Router, cfg *config.Components) {
	cfgAPI := &ConfigAPI{cfg: cfg}
	router.GET("/api/v1/config/reveal", cfgAPI.RevealConfig)
	router.GET("/api/v1/config/mermaid", cfgAPI.MermaidConfig)
}

func (a *ConfigAPI) RevealConfig(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(a.cfg.Reveal); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *ConfigAPI) MermaidConfig(writer http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	writer.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(writer)
	if err := enc.Encode(a.cfg.Mermaid); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
