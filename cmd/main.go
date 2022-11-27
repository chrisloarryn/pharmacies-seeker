package main

import (
	"pharmacies-seeker/internal/http/server"
	"pharmacies-seeker/internal/infraestucture/dependencies"

	"fmt"
)

func main() {
	fmt.Println("Running...")

	container := dependencies.NewContainer()

	server.Run(container)
}
