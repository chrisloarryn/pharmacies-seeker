## go command
GOCMD=go

server:
	${GOCMD} run cmd/main.go

test:
	${GOCMD} test -v ./...