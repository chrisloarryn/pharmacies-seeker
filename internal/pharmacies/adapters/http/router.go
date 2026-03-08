package pharmacyhttp

import "github.com/gofiber/fiber/v3"

func NewRouter(handler *Handler) *fiber.App {
	app := fiber.New()

	app.Get("/health/live", handler.Live)
	app.Get("/health/ready", handler.Ready)

	api := app.Group("/api/v1")
	api.Get("/pharmacies", handler.ListPharmacies)
	api.Get("/pharmacies/duty", handler.ListDutyPharmacies)
	api.Get("/pharmacies/:id", handler.GetPharmacyByID)

	return app
}
