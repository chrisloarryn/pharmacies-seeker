##
## Build
##
FROM golang:1.25-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /pharmacy-api cmd/main.go

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /pharmacy-api /pharmacy-api

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/pharmacy-api"]