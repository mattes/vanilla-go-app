build:
	-rm -r build
	GOOS=linux GOARCH=amd64 go build -a -o build/vanilla-go-app.linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -a -o build/vanilla-go-app.darwin-amd64 .

test:
	go test -v -race ./server 

run:
	go run main.go

build-docker:
	docker build -t templarbit/vanilla-go-app:latest .

push:
	docker push templarbit/vanilla-go-app:latest

.PHONY: build test run build-docker push
