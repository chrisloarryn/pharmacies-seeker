package handlers

import (
	"encoding/xml"
	"net/http"
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

func (handler *FindAllPharmaciesHandler) GetAllPharmacies(ctx *fiber.Ctx) error {
	communeName := ctx.Query("commune", "")
	responseType := ctx.Query("type")

	pharmacies, err := handler.uc.Execute(ctx.Context(), communeName)

	if responseType == "xml" {
		xmlBytes, err := xml.Marshal(pharmacies)
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
