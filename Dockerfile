FROM golang:1.10-alpine

WORKDIR $GOPATH/src/github.com/templarbit/vanilla-go-app

COPY . .

RUN go build

RUN go install

CMD ["vanilla-go-app"] 
