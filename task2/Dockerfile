FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd/reader/main.go ./

RUN go build -o /main

CMD [ "/main" ]