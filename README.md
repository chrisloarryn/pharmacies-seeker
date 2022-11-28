# pharmacies-seeker

## Description

This is a simple script to find pharmacies in a given city.

## Development

### Usage (go ~ docker ~ docker-compose)

```bash
## go executable
$ go run cmd/main.go
```

```bash
## Build the image
$ docker build -t pharmacies-seeker .

## Run the container with the image
$ docker run -it pharmacies-seeker
 ```

```bash
## Build the image and run the container (add '-d' for detached mode)
$ docker-compose up --build 
```

### Configs

#### Configuration file format in YAML, loaded from `./internal/shared/config/config.yml` file.

```yaml
server:
  port: 8080

api:
  pharmacy:
    url: https://farmanet.minsal.cl/index.php/ws/getLocales
```

### Tools

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Fiber](https://gofiber.io/)
- [Viper](https://github.com/spf13/viper)
- [Mock](github.com/golang/mock)
- [Testify](github.com/stretchr/testify)
### Requirements

- Go 1.18
- Docker
- Docker Compose
- Make
- Git

### Endpoints

#### host: `localhost:8080`
#### prefix: `/api/v1`

- `GET /pharmacies?commune={commune}&type={type}`

### Examples

##### Request for pharmacies in a given `commune` and type `json`

```bash
$ curl -X GET "http://localhost:8080/api/v1/pharmacies?commune=CONCON&type="
```

#### type=JSON (default)

```json5
{
  "message": "OK",
  "data": [
    {
      "local_nombre": "CRUZ VERDE",
      "comuna_nombre": "CONCON",
      "local_direccion": "LOS GINKOS 5 LOCAL 11,12,13",
      "local_telefono": "+56322857355"
    },
    {
      "local_nombre": "CRUZ VERDE",
      "comuna_nombre": "CONCON",
      "local_direccion": "AV. CON CON REÑACA 3850 LOCAL 1013",
      "local_telefono": "+56322858104"
    },
  ]
}
````

--

##### Request for pharmacies in a given `commune` and type `xml`

```bash
$ curl -X GET "http://localhost:8080/api/v1/pharmacies?commune=CONCON&type=xml"
```

#### type=XML

```xml

<Pharmacies>
    <Pharmacy>
        <local_nombre>CRUZ VERDE</local_nombre>
        <comuna_nombre>CONCON</comuna_nombre>
        <local_direccion>LOS GINKOS 5 LOCAL 11,12,13</local_direccion>
        <local_telefono>+56322857355</local_telefono>
    </Pharmacy>
    <Pharmacy>
        <local_nombre>CRUZ VERDE</local_nombre>
        <comuna_nombre>CONCON</comuna_nombre>
        <local_direccion>AV. CON CON REÑACA 3850 LOCAL 1013</local_direccion>
        <local_telefono>+56322858104</local_telefono>
    </Pharmacy>
</Pharmacies>
```

### Swagger

```shell

# To generate a swagger spec document for a go application
$ swagger generate spec -o ./swagger.json

# Spec validation tool
$ swagger validate https://raw.githubusercontent.com/swagger-api/swagger-spec/master/examples/v2.0/json/petstore-expanded.json

# Generate a client from a swagger spec
$ swagger generate client [-f ./swagger.json] -A [application-name [--principal [principal-name]]

# Generate a server from a swagger spec
$ swagger generate server -f ./swagger.json -A [application-name] [--principal [principal-name]]
```