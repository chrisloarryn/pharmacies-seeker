package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
)

type catalogRepositoryStub struct {
	lists map[domain.CatalogKind][]domain.Pharmacy
	byID  map[string]domain.Pharmacy
}

func (s *catalogRepositoryStub) Replace(context.Context, domain.CatalogKind, []domain.Pharmacy) error {
	return nil
}

func (s *catalogRepositoryStub) List(_ context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	pharmacies, ok := s.lists[kind]
	if !ok {
		return nil, domain.ErrCatalogUnavailable
	}
	return append([]domain.Pharmacy(nil), pharmacies...), nil
}

func (s *catalogRepositoryStub) GetByID(_ context.Context, kind domain.CatalogKind, id string) (domain.Pharmacy, error) {
	if kind != domain.CatalogRegular {
		return domain.Pharmacy{}, domain.ErrInvalidCatalogKind
	}
	pharmacy, ok := s.byID[id]
	if !ok {
		return domain.Pharmacy{}, domain.ErrPharmacyNotFound
	}
	return pharmacy, nil
}

func TestCatalogServiceListFiltersByQuery(t *testing.T) {
	repository := &catalogRepositoryStub{
		lists: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {
				{ID: "1", Name: "Cruz Verde", Commune: "Santiago", Address: "A", Phone: "1"},
				{ID: "2", Name: "Salcobrand", Commune: "Providencia", Address: "B", Phone: "2"},
			},
		},
	}

	service := app.NewCatalogService(repository)

	pharmacies, err := service.List(context.Background(), domain.CatalogRegular, domain.NewQuery("santi", "verde"))

	require.NoError(t, err)
	require.Len(t, pharmacies, 1)
	assert.Equal(t, "1", pharmacies[0].ID)
}

func TestCatalogServiceListRejectsInvalidCatalogKind(t *testing.T) {
	service := app.NewCatalogService(&catalogRepositoryStub{})

	_, err := service.List(context.Background(), domain.CatalogKind("bad"), domain.NewQuery("", ""))

	require.ErrorIs(t, err, domain.ErrInvalidCatalogKind)
}

func TestCatalogServiceListPropagatesRepositoryError(t *testing.T) {
	service := app.NewCatalogService(&catalogRepositoryStub{})

	_, err := service.List(context.Background(), domain.CatalogDuty, domain.NewQuery("", ""))

	require.ErrorIs(t, err, domain.ErrCatalogUnavailable)
}

func TestCatalogServiceGetByIDReturnsPharmacy(t *testing.T) {
	repository := &catalogRepositoryStub{
		byID: map[string]domain.Pharmacy{
			"abc": {ID: "abc", Name: "Cruz Verde", Commune: "Santiago"},
		},
	}

	service := app.NewCatalogService(repository)

	pharmacy, err := service.GetByID(context.Background(), " abc ")

	require.NoError(t, err)
	assert.Equal(t, "abc", pharmacy.ID)
}

func TestCatalogServiceGetByIDPropagatesNotFound(t *testing.T) {
	service := app.NewCatalogService(&catalogRepositoryStub{byID: map[string]domain.Pharmacy{}})

	_, err := service.GetByID(context.Background(), "missing")

	require.ErrorIs(t, err, domain.ErrPharmacyNotFound)
}

func TestCatalogServiceGetByIDRejectsInvalidID(t *testing.T) {
	service := app.NewCatalogService(&catalogRepositoryStub{})

	_, err := service.GetByID(context.Background(), " ")

	require.ErrorIs(t, err, domain.ErrInvalidID)
}
