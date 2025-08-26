package handlers

import (
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/infraestucture/dependencies"
)

type DutyPharmaciesHandler struct {
	fs fetcher.Service
}

func NewDutyPharmaciesHandler(container dependencies.Container) *DutyPharmaciesHandler {
	return &DutyPharmaciesHandler{fs: container.FetcherService()}
}

// GetDutyPharmacies
// swagger: route GET /pharmacies/24h getDutyPharmacies
//
// # Get 24h pharmacies
//
// Responses:
//   - 200: GetAllPharmaciesResponse
func (h *DutyPharmaciesHandler) GetDutyPharmacies(ctx *fiber.Ctx) error {
	commune := ctx.Query("commune", "")
	name := ctx.Query("name", "")
	responseType := ctx.Query("type")

	list, err := h.fs.RetrievePharmacies24h(ctx.Context())
	if err != nil {
		return reply(ctx, http.StatusInternalServerError, err.Error(), nil)
	}

	// filters
	if commune != "" {
		filtered := make([]pharmacy.Pharmacy, 0, len(list))
		needle := strings.ToLower(commune)
		for _, p := range list {
			if strings.Contains(strings.ToLower(p.ComunaNombre), needle) {
				filtered = append(filtered, p)
			}
		}
		list = filtered
	}
	if name != "" {
		needle := strings.ToLower(name)
		filtered := make([]pharmacy.Pharmacy, 0, len(list))
		for _, p := range list {
			if strings.Contains(strings.ToLower(p.LocalNombre), needle) {
				filtered = append(filtered, p)
			}
		}
		list = filtered
	}

	if responseType == "xml" {
		var pharmaciesXML pharmacy.Pharmacies
		pharmaciesXML.Pharmacies = list
		xmlBytes, err := xml.Marshal(pharmaciesXML)
		if err != nil {
			return err
		}
		ctx.Set("Content-Type", "application/xml")
		return ctx.SendString(string(xmlBytes))
	}

	return reply(ctx, http.StatusOK, "OK", list)
}
