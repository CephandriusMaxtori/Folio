package Folio

import "embed"

//go:embed web/index.html web/css web/js
var WebFiles embed.FS
