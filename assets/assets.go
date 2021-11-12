package assets

import "embed"

var (
	//go:embed template/reveal-markdown.tmpl
	Template []byte
	//go:embed web reveal
	Assets embed.FS
)
