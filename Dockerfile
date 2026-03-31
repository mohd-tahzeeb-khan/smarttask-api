# ── Build stage ──────────────────────────────────────────────────────────────
FROM golang:1.23-bookworm AS builder

RUN apt-get update && apt-get install -y \
    gcc \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=1
ENV GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o smarttask ./cmd/main.go

# ── Run stage ─────────────────────────────────────────────────────────────────
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libsqlite3-0 \
    ca-certificates \
    tzdata \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/smarttask .

RUN mkdir -p /app/data

EXPOSE 8080

ENV PORT=8080
ENV GIN_MODE=release
ENV DB_PATH=/app/data/smarttask.db

CMD ["./smarttask"]