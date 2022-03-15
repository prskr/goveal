package web

import "embed"

//go:embed js/* css/* index.gohtml
var WebFS embed.FS
