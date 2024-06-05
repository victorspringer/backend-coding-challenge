FROM golang:1.22.3-alpine

WORKDIR /app

COPY configs ./

COPY vendor ./

COPY . .

RUN go build -o ratingservice ./cmd/service

EXPOSE 8082

CMD ["./ratingservice"]
