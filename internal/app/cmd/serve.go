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

package cmd

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/baez90/go-reveal-slides/internal/app/rendering"
	"github.com/baez90/go-reveal-slides/internal/app/routing"
	"github.com/markbates/pkger"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	host     string
	port     uint16
	serveCmd = &cobra.Command{
		Use:   "serve",
		Args:  cobra.ExactArgs(1),
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			router := &routing.RegexpRouter{}
			var tmplRenderer rendering.RevealRenderer
			if tmplRenderer, err = rendering.NewRevealRenderer(&params); err != nil {
				log.Errorf("Failed to initialize reveal renderer due to error: %v", err)
				return
			}

			log.Info("Setup template renderer")
			if err = router.AddRule(`^(/(index.html(l)?)?)?$`, tmplRenderer); err != nil {
				return
			}

			var markdownHandler rendering.MarkdownHandler
			if markdownHandler, err = rendering.NewMarkdownHandler(args[0]); err != nil {
				log.Errorf("Failed to initialize reveal renderer due to error: %v", err)
				return
			}

			// single file handler that only delivers the single Markdown file containing the slides
			log.Info("Setup markdown handler for any *.md file...")
			if err = router.AddRule(".*\\.md$", markdownHandler); err != nil {
				return
			}

			layeredHandler := &routing.LayeredHandler{}
			layeredHandler.AddHandlers(pkger.Dir("/assets/reveal"), http.Dir("."))

			log.Info("Setup local file system and Reveal assets under / route")
			if err = router.AddRule("/.+", layeredHandler); err != nil {
				return
			}

			// start HTTP server
			hostPort := fmt.Sprintf("%s:%d", host, port)
			log.Infof("Running at addr http://%s/", hostPort)
			if err = http.ListenAndServe(hostPort, router); err != nil && errors.Is(err, http.ErrServerClosed) {
				log.Errorf("Error while running serve command: %v", err)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&host, "host", "localhost", "host the CLI should listen on")
	serveCmd.Flags().Uint16Var(&port, "port", 2233, "port the CLI should listen on")
}
