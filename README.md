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
### Requirements

- Go 1.18
- Docker
- Docker Compose
- Make
- Git