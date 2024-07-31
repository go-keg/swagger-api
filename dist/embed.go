package dist

import "embed"

//go:embed *.png *.css *.js
var SwagFS embed.FS
