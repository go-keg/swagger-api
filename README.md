# swagger-api


## Install

### protoc-gen-openapi
```shell
go install github.com/google/gnostic/cmd/protoc-gen-openapi
```

## QuickStart

```shell
go get github.com/go-keg/swagger-api
```

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
	mux.Handle("/swagger/", swagger.Handler(api.FS))
	err := http.ListenAndServe(":8088", mux)
	if err != nil {
		panic(err)
	}
}
```

SwaggerUI: http://127.0.0.1:8088/swagger/ui


### Kratos example
```go
package main

import (
	"github.com/go-keg/swagger-api"
	"github.com/go-keg/swagger-api/examples/apis/api"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func main() {
	srv := http.NewServer(http.Address(":8080"))
	// http://127.0.0.1:8080/swagger/ui
	srv.HandlePrefix("/swagger", swagger.Handler(api.FS))

	// http://127.0.0.1:8080/swag
	srv.HandlePrefix("/swag", swagger.Handler(
		api.FS,
		swagger.SetPrefix("/swag"),    // default: /swagger
		swagger.SetSwaggerUIPath("/"), // default: /ui
		swagger.SetURLs([]swagger.OpenapiURL{
			{Name: "账号服务", URL: "account-interface/v1/account.openapi.yaml"},
		}), // 手动指定服务名称和路径，默认为所有引入的文件
	))
	app := kratos.New(kratos.Name("example"), kratos.Server(srv))
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```