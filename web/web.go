package web

import "embed"

//go:embed js/* index.gohtml
var WebFS embed.FS
