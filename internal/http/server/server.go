package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"pharmacies-seeker/internal/http/server/handlers"
	"pharmacies-seeker/internal/infraestucture/dependencies"
	"pharmacies-seeker/internal/shared/constants"

	"github.com/gofiber/fiber/v2"
)

type ServerHTTP struct{}

func Run(container dependencies.Container) {
	r := fiber.New()

	// fill data
	list, err := container.FetcherService().RetrievePharmacies(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	container.PharmaciesRepository().LoadAll(context.TODO(), list)

	r.Get("/", pingpong)
	r.Get("/ping", pingpong)

	api := r.Group("/api")

	v1 := api.Group("/v1", func(c *fiber.Ctx) error { // middleware for /api/v1
		c.Set("Version", "v1")
		return c.Next()
	})

	getAllHandler := handlers.NewFindAllPharmaciesHandler(container)
	getOneHandler := handlers.NewFindOnePharmacyHandler(container)

	v1.Get("/pharmacies", getAllHandler.GetAllPharmacies)
	v1.Get("/pharmacies/commune", getOneHandler.FindOnePharmacy)

	port := os.Getenv(constants.Port)

	log.Fatal(r.Listen(fmt.Sprintf(":%s", port)))
}

func pingpong(ctx *fiber.Ctx) error {
	ctx.Status(http.StatusOK)
	err := ctx.JSON(
		fiber.Map{
			"message": "Pong",
		})
	if err != nil {
		return err
	}
	return nil
}
