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

package template

import (
	"fmt"
	"github.com/baez90/go-reveal-slides/internal/app/config"
	"github.com/gobuffalo/packr/v2"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

type RevealRenderer interface {
	http.Handler
	init() error
}

func NewRevealRenderer(markdownPath string, params *config.RevealParams) (renderer RevealRenderer, err error) {
	var info os.FileInfo
	info, err = os.Stat(markdownPath)
	if err != nil {
		return
	}

	if info.IsDir() || path.Ext(info.Name()) != ".md" {
		err = fmt.Errorf("path %s did not pass sanity checks for markdown files", markdownPath)
		return
	}

	renderer = &revealRenderer{
		markdownPath: markdownPath,
		params:       params,
	}
	err = renderer.init()

	return
}

type revealRenderer struct {
	markdownPath     string
	template         *template.Template
	renderedTemplate string
	params           *config.RevealParams
}

func (renderer *revealRenderer) init() (err error) {
	templateBox := packr.New("template", "./../../../assets/template")
	templateString, err := templateBox.FindString("reveal-markdown.tmpl")
	if err != nil {
		return
	}

	renderer.template, err = template.New("index").Parse(templateString)
	return
}

func (renderer *revealRenderer) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	if renderer.template == nil {
		writeErrorResponse(500, "template is not set - probably error during startup", response)
		return
	}

	markdownContent, err := ioutil.ReadFile(renderer.markdownPath)
	if err != nil {
		writeErrorResponse(500, "failed to read markdown content", response)
		log.Errorf("Failed to read markdown content: %v", err)
		return
	}

	err = renderer.template.Execute(response, struct {
		Reveal       config.RevealParams
		MarkdownBody string
	}{Reveal: *renderer.params, MarkdownBody: string(markdownContent)})

	if err != nil {
		writeErrorResponse(500, "Failed to render Markdown to template", response)
		log.Errorf("Failed to render Markdown template: %v", err)
	}
}

func writeErrorResponse(code int, msg string, response http.ResponseWriter) {
	response.WriteHeader(code)
	_, err := response.Write([]byte(msg))
	log.Errorf("Failed to write error reponse: %v", err)
}
