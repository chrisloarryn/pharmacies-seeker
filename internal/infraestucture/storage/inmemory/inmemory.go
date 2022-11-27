package inmemory

import (
	"context"
	"errors"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/http/client/fetcherService"
)

type Repository struct {
	list    []pharmacy.Pharmacy
	fetcher fetcher.Service
}

func (r *Repository) Find(ctx context.Context, commune string) ([]pharmacy.Pharmacy, error) {
	if commune != "" {
		var locals []pharmacy.Pharmacy
		for _, v := range r.list {
			if v.ComunaNombre == commune {
				locals = append(locals, v)
			}
		}
		return locals, nil
	}

	return r.list, nil
}

func (r *Repository) FindOne(ctx context.Context, commune string) (pharmacy.Pharmacy, error) {
	if commune == "" {
		return pharmacy.Pharmacy{}, errors.New("commune is required")
	}
	for _, v := range r.list {
		if v.ComunaNombre == commune {
			return v, nil
		}
	}
	return pharmacy.Pharmacy{}, nil
}

func (r *Repository) LoadAll(ctx context.Context, pharmacies []pharmacy.Pharmacy) error {
	r.list = pharmacies
	return nil
}

func NewInMemoryRepository() pharmacy.Repository {
	return &Repository{
		list:    []pharmacy.Pharmacy{},
		fetcher: fetcherService.NewFetcherService(),
	}
}
