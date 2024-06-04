FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

COPY configs ./

RUN go mod download

COPY . .

RUN go build -o authenticationservice ./cmd/service

EXPOSE 8084

CMD ["./authenticationservice"]
