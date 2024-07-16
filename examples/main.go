package main

import (
	"github.com/eiixy/swagger-api"
	_ "github.com/eiixy/swagger-api/examples/apis/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
)

func main() {
	apisFS, err := fs.NewWithNamespace("apis")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagger.Handler(apisFS, []swagger.OpenapiURL{
		{"Account Interface", "/account-interface/v1/account.openapi.yaml"},
		{"Auth Interface", "/auth-interface/v1/auth.openapi.yaml"},
	},
		swagger.SetPrefix("/swagger"),
		swagger.SetSwaggerUIPath("ui"),
		swagger.SetOpenapiPath("apis"),
	))
	err = http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
