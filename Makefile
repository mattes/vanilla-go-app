build:
	-rm -r build
	GOOS=linux GOARCH=amd64 go build -a -o build/vanilla-go-app.linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -a -o build/vanilla-go-app.darwin-amd64 .

test:
	go test -v -race ./server 

run:
	go run main.go

.PHONY: build test run
