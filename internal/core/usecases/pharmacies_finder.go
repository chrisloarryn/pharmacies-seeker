package usecases

import (
	"context"
	"pharmacies-seeker/internal/core/domain/pharmacy"
)

// FinderAllPharmacies is the use case than find all pharmacy
type FinderAllPharmacies struct {
	pharmaciesRepository pharmacy.Repository
}

func NewFinderAllPharmacies(repository pharmacy.Repository) *FinderAllPharmacies {
	return &FinderAllPharmacies{
		repository,
	}
}

// Execute finder in the repository of pharmacy
func (f *FinderAllPharmacies) Execute(ctx context.Context, commune string) ([]pharmacy.Pharmacy, error) {
	if commune != "" {
		if err := pharmacy.ValidatePharmacyCommune(commune); err != nil {
			return nil, err
		}
	}
	pharmaciesList, err := f.pharmaciesRepository.Find(ctx, commune)
	if err != nil {
		return nil, err
	}
	return pharmaciesList, nil
}
