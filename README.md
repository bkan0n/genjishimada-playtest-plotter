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

### Docker (from GHCR)

```bash
docker pull ghcr.io/genjishimada/playtest-plotter:latest
docker run -p 8080:8080 ghcr.io/genjishimada/playtest-plotter:latest
```

Or with docker-compose:

```yaml
services:
  playtest-plotter:
    image: ghcr.io/genjishimada/playtest-plotter:latest
    ports:
      - "8080:8080"
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

## Releases

Releases are automated via GitHub Actions. When a version tag is pushed:

1. Docker image is built and pushed to GHCR with version tags
2. Linux amd64 binary is compiled and attached to the release
3. GitHub Release is created with auto-generated release notes

### Creating a release

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Available tags

| Tag | Description |
|-----|-------------|
| `latest` | Most recent release |
| `1.0.0` | Specific version |
| `1.0` | Latest patch of minor version |
