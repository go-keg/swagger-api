package main

import (
	"github.com/eiixy/swagger-api"
	"github.com/eiixy/swagger-api/examples/apis/api"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagger.Handler(http.FS(api.FS), []swagger.OpenapiURL{
		{"Account Interface", "/account-interface/v1/account.openapi.yaml"},
		{"Auth Interface", "/auth-interface/v1/auth.openapi.yaml"},
	},
		swagger.SetPrefix("/swagger"),
		swagger.SetSwaggerUIPath("ui"),
		swagger.SetOpenapiPath("apis"),
	))
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
