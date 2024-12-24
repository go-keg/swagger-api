package main

import (
	"net/http"

	"github.com/go-keg/swagger-api"
	"github.com/go-keg/swagger-api/examples/apis/api"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagger.Handler(http.FS(api.FS), []swagger.OpenapiURL{
		{Name: "Account Interface", URL: "/account-interface/v1/account.openapi.yaml"},
		{Name: "Auth Interface", URL: "/auth-interface/v1/auth.openapi.yaml"},
	},
	//swagger.SetPrefix("/swagger"),
	//swagger.SetSwaggerUIPath("ui"),
	//swagger.SetOpenapiPath("apis"),
	))
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
