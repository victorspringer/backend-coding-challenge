FROM golang:1.22.3-alpine

WORKDIR /app

COPY configs ./

COPY vendor ./

COPY . .

RUN go build -o movieservice ./cmd/service

EXPOSE 8083

CMD ["./movieservice"]
