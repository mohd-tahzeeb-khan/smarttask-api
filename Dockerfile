# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev build-base

WORKDIR /app

COPY go.mod go.sum ./

ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o smarttask ./cmd/main.go

# ── Run stage ─────────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/smarttask .

RUN mkdir -p /app/data

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE=release
ENV DB_PATH=/app/data/smarttask.db

CMD ["./smarttask"]