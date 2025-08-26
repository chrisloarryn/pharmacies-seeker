package fetcher

import (
	"context"
	"pharmacies-seeker/internal/core/domain/pharmacy"
)

//go:generate mockgen -package mocks -destination mocks/pharmacy_service_mocks.go . Service

// Service is the service abstraction to Pharmacy
type Service interface {
	RetrievePharmacies(ctx context.Context) ([]pharmacy.Pharmacy, error)
	RetrievePharmacies24h(ctx context.Context) ([]pharmacy.Pharmacy, error)
}
