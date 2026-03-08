# pharmacies-seeker

API en Go para consultar farmacias y farmacias de turno usando un proveedor remoto, con cache en memoria y refresh periódico.

## Arquitectura

```text
cmd/api
internal/pharmacies/{domain,app,adapters}
internal/platform/{config,http,scheduler}
```

## Requisitos

- Go 1.26.1
- Docker
- Docker Compose
- Make

## Configuración

La app carga defaults desde `internal/platform/config/properties.yml` y permite override por variables de entorno.

Variables soportadas:

- `PORT`
- `PROVIDER_REGULAR_URL`
- `PROVIDER_DUTY_URL`
- `SYNC_INTERVAL`
- `SYNC_TIMEOUT`

Defaults:

- `sync.interval=15m`
- `sync.timeout=10s`

## Desarrollo

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

## Respuestas

Listado:

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

Detalle:

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

## Calidad

```bash
go test ./...
go vet ./...
go build ./...
```
