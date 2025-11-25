# ビルド用
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# 実行用
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

# DB_DSN は Docker Compose または docker run で渡す
ENV DB_DSN=""

EXPOSE 8081
CMD ["./main"]

