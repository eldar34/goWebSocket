FROM golang:1.16-alpine

WORKDIR /app

COPY ./ /app

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

RUN apk add acf-openssl

ENTRYPOINT CompileDaemon --build="go build" --command=./main