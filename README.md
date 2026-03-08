# pharmacies-seeker

Go API for querying pharmacies and duty pharmacies using a remote provider, with in-memory caching and periodic refresh.

## Architecture

```text
cmd/api
internal/pharmacies/{domain,app,adapters}
internal/platform/{config,http,scheduler}
```

## Requirements

- Go 1.26.1
- Docker
- Docker Compose
- Make

## Configuration

The application loads defaults from `internal/platform/config/properties.yml` and allows overrides through environment variables.

Supported variables:

- `PORT`
- `PROVIDER_REGULAR_URL`
- `PROVIDER_DUTY_URL`
- `SYNC_INTERVAL`
- `SYNC_TIMEOUT`

Defaults:

- `sync.interval=15m`
- `sync.timeout=10s`

## Development

```bash
make run
make test
make build
```

```bash
docker compose up --build
```

## Endpoints

- `GET /health/live`
- `GET /health/ready`
- `GET /api/v1/pharmacies?commune={commune}&name={name}`
- `GET /api/v1/pharmacies/{id}`
- `GET /api/v1/pharmacies/duty?commune={commune}&name={name}`

## Responses

List:

```json
{
  "data": [
    {
      "id": "1",
      "name": "CRUZ VERDE",
      "commune": "SANTIAGO",
      "address": "AV. PROVIDENCIA 100",
      "phone": "+56212345678"
    }
  ]
}
```

Detail:

```json
{
  "data": {
    "id": "1",
    "name": "CRUZ VERDE",
    "commune": "SANTIAGO",
    "address": "AV. PROVIDENCIA 100",
    "phone": "+56212345678"
  }
}
```

Error:

```json
{
  "error": {
    "code": "not_found",
    "message": "pharmacy not found"
  }
}
```

## Quality

```bash
go test ./...
go vet ./...
go build ./...
```
