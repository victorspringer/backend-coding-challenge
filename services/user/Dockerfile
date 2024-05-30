FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

COPY configs ./

RUN go mod download

COPY . .

RUN go build -o userservice ./cmd/service

EXPOSE 8081

CMD ["./userservice"]
