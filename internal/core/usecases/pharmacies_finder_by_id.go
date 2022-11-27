package usecases

import (
	"context"
	"errors"
	"pharmacies-seeker/internal/core/domain/pharmacy"
)

// FinderPharmaciesByCommune is the use case than find all pharmacies
type FinderPharmaciesByCommune struct {
	pharmaciesRepository pharmacy.Repository
}

func NewFinderPharmaciesByCommune(repository pharmacy.Repository) *FinderPharmaciesByCommune {
	return &FinderPharmaciesByCommune{
		repository,
	}
}

// Execute finder a pharmacy by his ID in the repository of pharmacies
func (pf *FinderPharmaciesByCommune) Execute(ctx context.Context, commune string) (pharmacy.Pharmacy, error) {
	if commune == "" {
		return pharmacy.Pharmacy{}, errors.New("commune is required")
	}
	if err := pharmacy.ValidatePharmacyCommune(commune); err != nil {
		return pharmacy.Pharmacy{}, err
	}
	result, err := pf.pharmaciesRepository.FindOne(ctx, commune)
	if err != nil {
		return pharmacy.Pharmacy{}, err
	}
	return result, nil
}
