# Playtest Plotter

A Go microservice that generates WebP chart images showing player votes on map difficulty levels.

## Prerequisites

### Local Development

Install Cairo and libwebp:

```bash
# macOS
brew install cairo webp

# Ubuntu/Debian
apt-get install libcairo2-dev libwebp-dev
```

### Docker

No prerequisites - dependencies are bundled in the image.

## Running

### Local

```bash
go build -o chart-service .
./chart-service
```

### Docker

```bash
docker build -t playtest-plotter .
docker run -p 8080:8080 playtest-plotter
```

## API

### POST /chart

Generate a difficulty vote chart.

**Request:**
```json
{
  "votes": {
    "Medium": 15,
    "Medium +": 25,
    "Hard -": 30,
    "Hard": 20,
    "Hard +": 10
  }
}
```

**Response:** `image/webp`

**Valid difficulty levels:**
`Easy -`, `Easy`, `Easy +`, `Medium -`, `Medium`, `Medium +`, `Hard -`, `Hard`, `Hard +`, `Very Hard -`, `Very Hard`, `Very Hard +`, `Extreme -`, `Extreme`, `Extreme +`, `Hell`

### GET /health

Health check endpoint. Returns `{"status": "ok"}`.

## Example

```bash
curl -X POST http://localhost:8080/chart \
  -H "Content-Type: application/json" \
  -d '{"votes":{"Medium":15,"Hard -":30,"Hard":20}}' \
  -o chart.webp
```
