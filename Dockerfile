# Build stage
FROM golang:1.23-bookworm AS builder

RUN apt-get update && apt-get install -y \
    libcairo2-dev \
    libwebp-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /chart-service .

# Runtime stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    libcairo2 \
    libwebp7 \
    fontconfig \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /chart-service /chart-service
COPY fonts/ /usr/local/share/fonts/
RUN fc-cache -f -v

EXPOSE 8080
CMD ["/chart-service"]
