package handlers

import (
	"encoding/xml"
	"net/http"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/core/usecases"
	"pharmacies-seeker/internal/infraestucture/dependencies"

	"github.com/gofiber/fiber/v2"
)

type FindAllPharmaciesHandler struct {
	uc *usecases.FinderAllPharmacies // .FinderAllPharmacies
}

func NewFindAllPharmaciesHandler(container dependencies.Container) *FindAllPharmaciesHandler {
	return &FindAllPharmaciesHandler{
		uc: usecases.NewFinderAllPharmacies(container.PharmaciesRepository()),
	}
}

// swagger: route GET /pharmacies getAllPharmacies
//
// # Get all pharmacies
//
// Responses:
//
// - 200: GetAllPharmaciesResponse
func (handler *FindAllPharmaciesHandler) GetAllPharmacies(ctx *fiber.Ctx) error {
	communeName := ctx.Query("commune", "")
	responseType := ctx.Query("type")

	pharmacies, err := handler.uc.Execute(ctx.Context(), communeName)

	if responseType == "xml" {
		var pharmaciesXML pharmacy.Pharmacies

		pharmaciesXML.Pharmacies = pharmacies

		xmlBytes, err := xml.Marshal(pharmaciesXML)
		if err != nil {
			return err
		}
		ctx.Set("Content-Type", "application/xml")
		return ctx.SendString(string(xmlBytes))
	}

	if err != nil {
		return reply(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	return reply(ctx, http.StatusOK, "OK", pharmacies)
}
