FROM golang:1.10-alpine

RUN apk --no-cache add curl

WORKDIR $GOPATH/src/github.com/templarbit/vanilla-go-app

COPY . .

RUN go build

RUN go install

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=2 \
  CMD curl -sf http://127.0.0.1:8080 || exit 1

CMD ["vanilla-go-app"] 
