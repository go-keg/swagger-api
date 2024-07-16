# swagger-api


## Install

### protoc-gen-openapi
```shell
go install github.com/google/gnostic/cmd/protoc-gen-openapi
```

### statik
```shell
go get github.com/rakyll/statik
go install github.com/rakyll/statik
```

## QuickStart

[examples/apis/Makefile](./examples/apis/Makefile)
```shell
protoc --proto_path=./api/ \
    --proto_path=./third_party \
    --openapi_out=output_mode=source_relative:./api/ \
    $(API_PROTO_FILES)
statik -src=./api -include=*.openapi.yaml -ns apis	
```

[examples/main.go](./examples/main.go)
```go
package main

import (
	"github.com/eiixy/swagger-api"
	_ "github.com/eiixy/swagger-api/examples/statik"
	"github.com/rakyll/statik/fs"
	"net/http"
)

func main() {
	apisFS, err := fs.NewWithNamespace("apis")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", swagger.Handler(apisFS, []swagger.OpenapiURL{
		{"Account Interface", "/account-interface/v1/account.openapi.yaml"},
		{"Auth Interface", "/auth-interface/v1/auth.openapi.yaml"},
	}))
	err = http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
```