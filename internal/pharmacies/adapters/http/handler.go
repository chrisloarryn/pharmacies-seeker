package pharmacyhttp

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v3"

	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
)

type Handler struct {
	catalogs *app.CatalogService
	sync     *app.SyncService
}

func NewHandler(catalogs *app.CatalogService, sync *app.SyncService) *Handler {
	return &Handler{
		catalogs: catalogs,
		sync:     sync,
	}
}

func (h *Handler) Live(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(healthResponse{Status: "live"})
}

func (h *Handler) Ready(ctx fiber.Ctx) error {
	if !h.sync.Ready() {
		return writeError(ctx, http.StatusServiceUnavailable, "not_ready", domain.ErrCatalogUnavailable.Error())
	}

	return ctx.Status(http.StatusOK).JSON(healthResponse{Status: "ready"})
}

func (h *Handler) ListPharmacies(ctx fiber.Ctx) error {
	if !h.sync.Ready() {
		return writeError(ctx, http.StatusServiceUnavailable, "not_ready", domain.ErrCatalogUnavailable.Error())
	}

	pharmacies, err := h.catalogs.List(
		ctx.Context(),
		domain.CatalogRegular,
		domain.NewQuery(ctx.Query("commune", ""), ctx.Query("name", "")),
	)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(newPharmacyListResponse(pharmacies))
}

func (h *Handler) GetPharmacyByID(ctx fiber.Ctx) error {
	if !h.sync.Ready() {
		return writeError(ctx, http.StatusServiceUnavailable, "not_ready", domain.ErrCatalogUnavailable.Error())
	}

	pharmacy, err := h.catalogs.GetByID(ctx.Context(), ctx.Params("id"))
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(newPharmacyResponse(pharmacy))
}

func (h *Handler) ListDutyPharmacies(ctx fiber.Ctx) error {
	if !h.sync.Ready() {
		return writeError(ctx, http.StatusServiceUnavailable, "not_ready", domain.ErrCatalogUnavailable.Error())
	}

	pharmacies, err := h.catalogs.List(
		ctx.Context(),
		domain.CatalogDuty,
		domain.NewQuery(ctx.Query("commune", ""), ctx.Query("name", "")),
	)
	if err != nil {
		return handleError(ctx, err)
	}

	return ctx.Status(http.StatusOK).JSON(newPharmacyListResponse(pharmacies))
}

func handleError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidCatalogKind), errors.Is(err, domain.ErrInvalidID):
		return writeError(ctx, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, domain.ErrPharmacyNotFound):
		return writeError(ctx, http.StatusNotFound, "not_found", err.Error())
	case errors.Is(err, domain.ErrCatalogUnavailable):
		return writeError(ctx, http.StatusServiceUnavailable, "not_ready", err.Error())
	default:
		return writeError(ctx, http.StatusInternalServerError, "internal_error", "internal server error")
	}
}

func writeError(ctx fiber.Ctx, statusCode int, code, message string) error {
	return ctx.Status(statusCode).JSON(errorEnvelope{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}
