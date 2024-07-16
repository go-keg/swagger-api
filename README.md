# swagger-api


## 安装相关依赖

### protoc-gen-openapi
```shell
go install github.com/google/gnostic/cmd/protoc-gen-openapi
```

### statik
```shell
go get github.com/rakyll/statik
go install github.com/rakyll/statik
```

### [example](./examples/main.go)

```go
//go:generate statik -src=./apis/api -include=*.openapi.yaml -ns apis
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