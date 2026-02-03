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

Releases are automated via [Release Please](https://github.com/googleapis/release-please) and GitHub Actions.

### How it works

1. Push commits to `main` using [conventional commits](https://www.conventionalcommits.org/)
2. Release Please automatically creates/updates a release PR
3. Merge the PR to trigger a release
4. Docker image + binary are built and published

### Conventional commits

| Prefix | Version bump | Example |
|--------|--------------|---------|
| `fix:` | Patch (1.0.0 → 1.0.1) | `fix: handle empty votes` |
| `feat:` | Minor (1.0.0 → 1.1.0) | `feat: add PNG output` |
| `feat!:` | Major (1.0.0 → 2.0.0) | `feat!: change API response format` |
| `chore:` | No release | `chore: update deps` |

### Docker tags

| Tag | Description |
|-----|-------------|
| `latest` | Most recent release |
| `1.0.0` | Specific version |
| `1.0` | Latest patch of minor version |
| `main` | Latest commit on main branch |
