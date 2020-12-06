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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/baez90/goveal/internal/app/rendering"
	"github.com/fsnotify/fsnotify"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	workingDir string
	rootCmd    = &cobra.Command{
		Use:   "goveal",
		Short: "goveal is a small reveal.js server",
		Long: `goveal is a single static binary to host your reveal.js based markdown presentation.
It is running a small web server that loads your markdown file, renders a complete HTML page and delivers it including all the reveal.js assets.
It is not required to restart the server when you edit the markdown - a simple reload of the page is doing all the required magic.`,
	}
	params rendering.RevealParams
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)

	var err error
	workingDir, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-reveal-slides.yaml)")
	rootCmd.PersistentFlags().StringVar(&workingDir, "working-dir", workingDir, "working directory to use")
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	var err error
	if workingDir, err = filepath.Abs(workingDir); err != nil {
		log.Warnf("Failed to determine absolute path for working dir %s: %v", workingDir, err)
		return
	}

	var cwd string
	if cwd, err = os.Getwd(); err != nil {
		log.Warnf("Failed to determine current working directory")
		return
	}

	if cwd != workingDir {
		if err = os.Chdir(workingDir); err != nil {
			log.Warnf("Failed to change working directory to %s", workingDir)
		}
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Infof("Failed to determine home directory: %v", err)
		} else {
			viper.AddConfigPath(home)
		}

		viper.AddConfigPath(workingDir)
		viper.SetConfigName("goveal")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())

		log.Info("Starting to watch config file...")
		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			log.Info("Noticed configuration change...")
			params.Load()
		})
	}

	params.WorkingDirectory = workingDir
	params.Load()
}
