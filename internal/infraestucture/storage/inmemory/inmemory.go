package inmemory

import (
	"context"
	"errors"
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/http/client/fetcherService"
	"strings"
)

type Repository struct {
	list    []pharmacy.Pharmacy
	fetcher fetcher.Service
}

func (r *Repository) Find(ctx context.Context, commune string) ([]pharmacy.Pharmacy, error) {
	if commune != "" {
		var locals []pharmacy.Pharmacy
		needle := strings.ToLower(commune)
		for _, v := range r.list {
			if strings.Contains(strings.ToLower(v.ComunaNombre), needle) {
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
	needle := strings.ToLower(commune)
	for _, v := range r.list {
		if strings.Contains(strings.ToLower(v.ComunaNombre), needle) {
			return v, nil
		}
	}
	return pharmacy.Pharmacy{}, nil
}

func (r *Repository) LoadAll(ctx context.Context, pharmacies []pharmacy.Pharmacy) error {
	r.list = pharmacies
	return nil
}

func NewInMemoryRepository(c config.Config) pharmacy.Repository {
	return &Repository{
		list:    []pharmacy.Pharmacy{},
		fetcher: fetcherService.NewFetcherService(c),
	}
}
