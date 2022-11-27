package handlers

import (
	"fmt"
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

func (handler *FindAllPharmaciesHandler) GetAllPharmacies(ctx *fiber.Ctx) error {
	communeName := ctx.Query("commune", "")
	responseType := ctx.Query("type")

	pharmacies, err := handler.uc.Execute(ctx.Context(), communeName)
	if err != nil {
		ctx.JSON(err.Error())
		ctx.SendStatus(http.StatusInternalServerError)
		return err
	}

	if responseType == "xml" {
		fmt.Println("XML was not implemented yet")
		var xmlPharmacies []pharmacy.PharmacyXML
		for _, pharmacy := range pharmacies {
			//parse to xml
			r, err := pharmacy.ToXMLInterface()
			if err != nil {
				ctx.JSON(err.Error())
				ctx.SendStatus(http.StatusInternalServerError)
				return err
			}
			xmlPharmacies = append(xmlPharmacies, r)
		}
		// set content type to xml
		ctx.Set("Content-Type", "application/xml")
		ctx.JSON(xmlPharmacies)
		// convert to xml

		// return xml
		ctx.Format(pharmacies)
	}

	ctx.JSON(pharmacies)
	ctx.SendStatus(http.StatusOK)
	return nil
}
