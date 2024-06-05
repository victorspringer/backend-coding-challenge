FROM golang:1.22.3-alpine

WORKDIR /app

COPY configs ./

COPY vendor ./

COPY . .

RUN go build -o authenticationservice ./cmd/service

EXPOSE 8084

CMD ["./authenticationservice"]
