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
	"fmt"
	"github.com/baez90/go-reveal-slides/internal/app/template"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"net/http"
	"os"

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
		Run: func(cmd *cobra.Command, args []string) {

			tmplRenderer, err := template.NewRevealRenderer(args[0], &params)

			if err != nil {
				log.Errorf("Failed to initialize reveal renderer due to error: %v", err)
				os.Exit(1)
			}

			log.Info("Setup reveal assets under route /reveal/ route...")
			revealBox := packr.New("reveal-assets", "./../../../assets/reveal")
			http.Handle("/reveal/", http.StripPrefix("/reveal/", http.FileServer(revealBox)))

			//
			log.Info("Setup static file serving under /local/ route...")
			fs := http.FileServer(http.Dir("."))
			http.Handle("/local/", http.StripPrefix("/local/", fs))

			http.Handle("/", tmplRenderer)

			hostPort := fmt.Sprintf("%s:%d", host, port)

			log.Infof("Running at addr http://%s/", hostPort)
			if err := http.ListenAndServe(hostPort, nil); err != nil {
				log.Error("Error while running serve command: %v", err)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&host, "host", "localhost", "host the CLI should listen on")
	serveCmd.Flags().Uint16Var(&port, "port", 2233, "port the CLI should listen on")
}
