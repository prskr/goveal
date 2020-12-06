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
	"os/exec"
	"runtime"

	"github.com/baez90/goveal/internal/app/server"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	host        string
	port        uint16
	openBrowser bool
	serveCmd    = &cobra.Command{
		Use:   "serve",
		Args:  cobra.ExactArgs(1),
		Short: "",
		Long:  ``,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var srv *server.HTTPServer
			srv, err = server.NewHTTPServer(server.Config{
				Host:         host,
				Port:         port,
				MarkdownPath: args[0],
				RevealParams: &params,
			})

			if err != nil {
				log.Errorf("Error while setting up server: %v", err)
				return
			}

			listenUrl := fmt.Sprintf("http://%s/", srv.ListenAddress())
			log.Infof("Going to listen on %s", listenUrl)

			if openBrowser {
				log.Info("Opening browser...")
				openBrowserInBackground(listenUrl)
			}

			if err = srv.Serve(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	serveCmd.Flags().BoolVar(&openBrowser, "open-browser", true, "if the browser should be opened at the URL")
}

func openBrowserInBackground(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Warn(err)
	}
}
