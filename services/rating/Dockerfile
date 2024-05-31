FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./

COPY configs ./

RUN go mod download

COPY . .

RUN go build -o ratingservice ./cmd/service

EXPOSE 8082

CMD ["./ratingservice"]
