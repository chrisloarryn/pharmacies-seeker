package main

import (
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/http/server"
	"pharmacies-seeker/internal/infraestucture/dependencies"

	"fmt"
)

// Pharmacies seeker api
//
// This is a simple api to search for pharmacies in Chile by commune.
//
//	Schemes: http
//	Host: localhost:8080
//	BasePath: /api/v1
//	Version: 0.0.1
//	License: MIT http://opensource.org/licenses/MIT
//	Contact: Cristobal Contreras <chrisloarryn@gmail.com>
//
//	Consumes:
//
//	Produces:
//	- application/json
//	- application/xml
//
// swagger:meta
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
