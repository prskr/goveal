// Copyright © 2019 Peter Kurfer peter.kurfer@googlemail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rendering

import (
	"html/template"
	"net/http"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"

	"github.com/baez90/goveal/assets"
)

type RevealRenderer interface {
	http.Handler
	init() error
}

func NewRevealRenderer(params *RevealParams) (renderer RevealRenderer, err error) {
	renderer = &revealRenderer{
		params: params,
	}
	err = renderer.init()

	return
}

type revealRenderer struct {
	template *template.Template
	params   *RevealParams
}

func (renderer *revealRenderer) init() (err error) {
	renderer.template, err = template.New("index").Funcs(sprig.FuncMap()).Parse(string(assets.Template))
	return
}

func (renderer *revealRenderer) ServeHTTP(response http.ResponseWriter, _ *http.Request) {
	if renderer.template == nil {
		writeErrorResponse(500, "rendering is not set - probably error during startup", response)
		return
	}

	err := renderer.template.Execute(response, struct {
		Reveal RevealParams
	}{Reveal: *renderer.params})
	if err != nil {
		writeErrorResponse(500, "Failed to render Markdown to rendering", response)
		log.Errorf("Failed to render Markdown rendering: %v", err)
	}
}

func writeErrorResponse(code int, msg string, response http.ResponseWriter) {
	response.WriteHeader(code)
	_, err := response.Write([]byte(msg))
	log.Errorf("Failed to write error reponse: %v", err)
}
