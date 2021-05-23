package static

import "embed"

// Assets represents the embedded files.
//go:embed *.tmpl pages/*.tmpl
var Assets embed.FS
