package main

import (
	"encoding/hex"
	"hash/fnv"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/baez90/goveal/api"
	"github.com/baez90/goveal/config"
	"github.com/baez90/goveal/events"
	"github.com/baez90/goveal/fs"
	"github.com/baez90/goveal/web"
)

var (
	workingDir string
	cfg        *config.Components
	serveCmd   = &cobra.Command{
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

			app := fiber.New(fiber.Config{
				Views: html.NewFileSystem(http.FS(web.WebFS), ".gohtml").
					AddFunc("fileId", func(fileName string) string {
						h := fnv.New32a()
						return hex.EncodeToString(h.Sum([]byte(path.Base(fileName))))
					}),
			})
			hub := events.NewEventHub(
				wdfs,
				fnv.New32a(),
				events.MutationReloadForFile(args[0]),
				events.MutationConfigReloadForFile(filepath.Base(cfg.ConfigFileInUse)),
			)

			api.NoCache(app)
			api.RegisterViews(app, wdfs, args[0], cfg)
			api.RegisterEventsAPI(app, hub, log.StandardLogger())
			api.RegisterConfigAPI(app, cfg)
			if err := api.RegisterStaticFileHandling(app, wdfs); err != nil {
				return err
			}

			return app.Listen(":3000")
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)
}
