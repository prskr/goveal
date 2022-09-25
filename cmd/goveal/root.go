// Copyright © 2019 Peter Kurfer
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

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"code.icb4dc0.de/prskr/goveal/config"
)

//nolint:lll // explanations are rather long
var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "goveal",
		Short: "goveal is a small reveal.js server",
		Long: `goveal is a single static binary to host your reveal.js based markdown presentation.
It is running a small web server that loads your markdown file, renders a complete HTML page and delivers it including all the reveal.js assets.
It is not required to restart the server when you edit the markdown - a simple reload of the page is doing all the required magic.`,
		PersistentPreRunE: func(*cobra.Command, []string) (err error) {
			log.SetFormatter(&log.TextFormatter{
				ForceColors: true,
			})

			if workingDir, err = os.Getwd(); err != nil {
				return err
			}

			cfg, err = config.Load(workingDir, cfgFile)
			return err
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-reveal-slides.yaml)")
}
