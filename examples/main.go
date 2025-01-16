package main

import (
	"net/http"

	"github.com/go-keg/swagger-api"
	"github.com/go-keg/swagger-api/examples/apis/api"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/swagger/", swagger.Handler(api.FS))
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
