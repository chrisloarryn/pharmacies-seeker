package app

import (
	"context"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type CatalogService struct {
	repository Repository
}

func NewCatalogService(repository Repository) *CatalogService {
	return &CatalogService{repository: repository}
}

func (s *CatalogService) List(ctx context.Context, kind domain.CatalogKind, query domain.Query) ([]domain.Pharmacy, error) {
	if err := kind.Validate(); err != nil {
		return nil, err
	}

	pharmacies, err := s.repository.List(ctx, kind)
	if err != nil {
		return nil, err
	}

	filtered := make([]domain.Pharmacy, 0, len(pharmacies))
	for _, pharmacy := range pharmacies {
		if domain.Matches(pharmacy, query) {
			filtered = append(filtered, pharmacy)
		}
	}

	return filtered, nil
}

func (s *CatalogService) GetByID(ctx context.Context, id string) (domain.Pharmacy, error) {
	normalizedID, err := domain.NormalizeID(id)
	if err != nil {
		return domain.Pharmacy{}, err
	}

	return s.repository.GetByID(ctx, domain.CatalogRegular, normalizedID)
}
