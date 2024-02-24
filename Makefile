build:
	@go build -o bin/gopictureuploader

run: build
	@./bin/gopictureuploader 

test: 
	@go test -v ./...
