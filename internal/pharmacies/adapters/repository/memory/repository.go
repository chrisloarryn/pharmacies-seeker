package memory

import (
	"context"
	"sync"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type Repository struct {
	mu       sync.RWMutex
	catalogs map[domain.CatalogKind]catalogSnapshot
}

type catalogSnapshot struct {
	items []domain.Pharmacy
	byID  map[string]domain.Pharmacy
}

func NewRepository() *Repository {
	return &Repository{
		catalogs: make(map[domain.CatalogKind]catalogSnapshot, 2),
	}
}

func (r *Repository) Replace(_ context.Context, kind domain.CatalogKind, pharmacies []domain.Pharmacy) error {
	if err := kind.Validate(); err != nil {
		return err
	}

	items := append([]domain.Pharmacy(nil), pharmacies...)
	byID := make(map[string]domain.Pharmacy, len(items))
	for _, pharmacy := range items {
		if pharmacy.ID != "" {
			byID[pharmacy.ID] = pharmacy
		}
	}

	r.mu.Lock()
	r.catalogs[kind] = catalogSnapshot{
		items: items,
		byID:  byID,
	}
	r.mu.Unlock()

	return nil
}

func (r *Repository) List(_ context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	if err := kind.Validate(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	snapshot, ok := r.catalogs[kind]
	r.mu.RUnlock()
	if !ok {
		return nil, domain.ErrCatalogUnavailable
	}

	return append([]domain.Pharmacy(nil), snapshot.items...), nil
}

func (r *Repository) GetByID(_ context.Context, kind domain.CatalogKind, id string) (domain.Pharmacy, error) {
	if err := kind.Validate(); err != nil {
		return domain.Pharmacy{}, err
	}

	normalizedID, err := domain.NormalizeID(id)
	if err != nil {
		return domain.Pharmacy{}, err
	}

	r.mu.RLock()
	snapshot, ok := r.catalogs[kind]
	r.mu.RUnlock()
	if !ok {
		return domain.Pharmacy{}, domain.ErrCatalogUnavailable
	}

	pharmacy, found := snapshot.byID[normalizedID]
	if !found {
		return domain.Pharmacy{}, domain.ErrPharmacyNotFound
	}

	return pharmacy, nil
}
