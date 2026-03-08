package app

import (
	"context"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type Provider interface {
	Fetch(ctx context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error)
}

type Repository interface {
	Replace(ctx context.Context, kind domain.CatalogKind, pharmacies []domain.Pharmacy) error
	List(ctx context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error)
	GetByID(ctx context.Context, kind domain.CatalogKind, id string) (domain.Pharmacy, error)
}
