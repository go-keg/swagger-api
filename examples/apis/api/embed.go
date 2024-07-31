package api

import "embed"

//go:embed account-interface/v1/*.openapi.yaml auth-interface/v1/*.openapi.yaml
var OpenapiFS embed.FS
