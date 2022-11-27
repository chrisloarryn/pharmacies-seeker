package main

import (
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/http/server"
	"pharmacies-seeker/internal/infraestucture/dependencies"

	"fmt"
)

func main() {
	fmt.Println("Running...")

	conf, err := config.LoadConfig("./internal/shared/config")
	if err != nil {
		fmt.Println("Error loading config: ", err)
		panic(err.Error())
	}

	fmt.Println("Config loaded", conf)

	container := dependencies.NewContainer(conf)

	server.Run(container)
}
