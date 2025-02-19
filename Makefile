clean: 
	@rm -rf bin/
build: clean
	@go build -o bin/md-to-blog
test: 
	@go test -v ./...
