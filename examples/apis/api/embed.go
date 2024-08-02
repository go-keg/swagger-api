package api

import "embed"

//go:embed account-interface/v1/account.openapi.yaml
//go:embed auth-interface/v1/*.openapi.yaml
var FS embed.FS
