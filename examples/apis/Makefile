API_PROTO_FILES=$(shell find api -name *.proto)

gen:
	protoc --proto_path=./api/ \
		--proto_path=./third_party \
		--openapi_out=output_mode=source_relative:./api/ \
		$(API_PROTO_FILES)
	statik -src=./api -include=*.openapi.yaml -ns apis

install:
	go install github.com/google/gnostic/cmd/protoc-gen-openapi