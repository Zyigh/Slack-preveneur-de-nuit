FROM golang:1.19 AS builder

WORKDIR /pierre
ADD / .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix nocgo -o ./pierre

CMD ./pierre
