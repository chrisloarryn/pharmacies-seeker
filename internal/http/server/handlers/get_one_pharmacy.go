package handlers

import (
	"net/http"
	"pharmacies-seeker/internal/core/usecases"
	"pharmacies-seeker/internal/infraestucture/dependencies"

	"github.com/gofiber/fiber/v2"
)

type FindOnePharmacyHandler struct {
	uc *usecases.FinderPharmaciesByCommune
}

func NewFindOnePharmacyHandler(container dependencies.Container) *FindOnePharmacyHandler {
	return &FindOnePharmacyHandler{
		uc: usecases.NewFinderPharmaciesByCommune(container.PharmaciesRepository()),
	}
}

func (handler *FindOnePharmacyHandler) FindOnePharmacy(ctx *fiber.Ctx) error {
	// communeName := ctx.Params("commune")
	communeName := ctx.Query("name")

	// TODO: Validate communeName
	pharmacy, err := handler.uc.Execute(ctx.Context(), communeName)
	if err != nil {
		ctx.JSON(err.Error())
		ctx.SendStatus(http.StatusInternalServerError)
		return err
	}
	ctx.JSON(pharmacy)
	ctx.SendStatus(http.StatusOK)
	return nil
}
