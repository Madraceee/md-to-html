clean: 
	@rm -rf bin/
build: clean
	@go build -o bin/md-to-html
test: 
	@go test -v ./...
test-coverage:
	@go test -cover -v ./...
build-wasm:
	@GOOS=js GOARCH=wasm go build -o site/main.wasm
generate-ast:
	@go build -o ./utils/tool ./utils/ && ./utils/tool
