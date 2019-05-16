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
	"github.com/baez90/go-reveal-slides/internal/app/config"
	"github.com/fsnotify/fsnotify"
	"os"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultTheme string = "white"
)

var (
	cfgFile             string
	theme               string
	codeTheme           string
	transition          string
	navigationMode      string
	horizontalSeparator string
	verticalSeparator   string
	rootCmd             = &cobra.Command{
		Use:   "goveal",
		Short: "goveal is a small reveal.js server",
		Long: `goveal is a single static binary to host your reveal.js based markdown presentation.
It is running a small web server that loads your markdown file, renders a complete HTML page and delivers it including all the reveal.js assets.
It is not required to restart the server when you edit the markdown - a simple reload of the page is doing all the required magic.`,
	}
	params config.RevealParams
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
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initLogging)

	rootCmd.PersistentFlags().StringVar(&theme, "theme", defaultTheme, "reveal.js theme to use")
	rootCmd.PersistentFlags().StringVar(&codeTheme, "code-theme", "monokai", "name of the code theme to use for highlighting")
	rootCmd.PersistentFlags().StringVar(&transition, "transition", "none", "transition effect to use")
	rootCmd.PersistentFlags().StringVar(&navigationMode, "navigationMode", "default", "determine the navigation mode to use ['default', 'linear', 'grid']")
	rootCmd.PersistentFlags().StringVar(&horizontalSeparator, "horizontal-separator", "===", "horizontal separator in slides")
	rootCmd.PersistentFlags().StringVar(&verticalSeparator, "vertical-separator", "---", "vertical separator in slides")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-reveal-slides.yaml)")

}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

		cwd, err := os.Getwd()
		if err != nil {
			log.Infof("Failed to determine current working directory: %v", err)
		} else {
			viper.AddConfigPath(cwd)
		}

		viper.SetConfigName("goveal")
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

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Errorf("Failed to bind flags to viper")
	}

	params.Load()

}
