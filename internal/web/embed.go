package web

import "embed"

// Pages embeds the four console pages served by the platform.
//
//go:embed *.html
var Pages embed.FS
