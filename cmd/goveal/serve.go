package main

import (
	"errors"
	"fmt"
	"hash/fnv"
	"net/http"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/baez90/goveal/api"
	"github.com/baez90/goveal/config"
	"github.com/baez90/goveal/events"
	"github.com/baez90/goveal/fs"
)

const (
	defaultListeningPort uint16 = 2233
	defaultHost                 = "127.0.0.1"
)

var (
	workingDir  string
	cfg         *config.Components
	host        string
	port        uint16
	openBrowser bool
	serveCmd    = &cobra.Command{
		Use:  "serve",
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) (err error) {
			var wdfs *fs.Watching

			if wdfs, err = fs.NewWatching(fs.Dir(workingDir)); err != nil {
				return err
			}

			defer multierr.AppendInvoke(&err, multierr.Close(wdfs))

			var mdFile fs.File
			if mdFile, err = wdfs.Open(args[0]); err != nil {
				return err
			}
			_ = mdFile.Close()

			router := httprouter.New()

			hub := events.NewEventHub(
				log.StandardLogger(),
				wdfs,
				fnv.New32a(),
				events.MutationReloadForFile(args[0]),
				events.MutationConfigReloadForFile(filepath.Base(cfg.ConfigFileInUse)),
			)

			api.RegisterViews(router, log.StandardLogger(), wdfs, args[0], cfg)
			api.RegisterEventsAPI(router, hub, log.StandardLogger())
			api.RegisterConfigAPI(router, cfg)

			handler := api.FileSystemMiddleware(router, wdfs)
			handler = api.NoCache(handler)

			if openBrowser {
				log.Info("Opening browser...")
				openBrowserInBackground(fmt.Sprintf("http://%s:%d", host, port))
			}

			log.Infof("Listening on %s:%d", host, port)
			if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), handler); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					return nil
				}
				return err
			}
			return nil
		},
	}
)

func init() {
	serveCmd.Flags().Uint16Var(&port, "port", defaultListeningPort, "port to listen on")
	serveCmd.Flags().StringVar(&host, "host", defaultHost, "address/hostname to bind on")
	serveCmd.Flags().BoolVar(&openBrowser, "open-browser", true, "Open browser when starting")
	rootCmd.AddCommand(serveCmd)
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
