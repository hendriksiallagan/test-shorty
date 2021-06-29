FROM golang:1.15-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum  ./

RUN go mod download

COPY . .

CMD ["go", "run", "main.go"]