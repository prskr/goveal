package assets

import "embed"

var (
	//go:embed reveal mermaid/mermaid.min.js
	Assets embed.FS
)
