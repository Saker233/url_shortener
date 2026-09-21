FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o url-shortener .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/url-shortener .

EXPOSE 8000

CMD ["./url-shortener"]