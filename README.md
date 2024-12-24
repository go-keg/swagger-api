# swagger-api


## Install

### protoc-gen-openapi
```shell
go install github.com/google/gnostic/cmd/protoc-gen-openapi
```

## QuickStart

[examples/apis/Makefile](./examples/apis/Makefile)
```shell
protoc --proto_path=./api/ \
    --proto_path=./third_party \
    --openapi_out=output_mode=source_relative:./api/ \
    $(API_PROTO_FILES)
```

[examples/main.go](./examples/main.go)
```go
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
	}))
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
```

SwaggerUI: http://127.0.0.1:8088/swagger/ui
