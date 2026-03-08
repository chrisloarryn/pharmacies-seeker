package memory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/domain"
)

func TestRepositoryReplaceAndList(t *testing.T) {
	repository := memory.NewRepository()

	err := repository.Replace(context.Background(), domain.CatalogRegular, []domain.Pharmacy{
		{ID: "1", Name: "Cruz Verde"},
	})

	require.NoError(t, err)

	pharmacies, err := repository.List(context.Background(), domain.CatalogRegular)
	require.NoError(t, err)
	require.Len(t, pharmacies, 1)
	assert.Equal(t, "1", pharmacies[0].ID)
}

func TestRepositoryListRejectsInvalidCatalog(t *testing.T) {
	repository := memory.NewRepository()

	_, err := repository.List(context.Background(), domain.CatalogKind("bad"))

	require.ErrorIs(t, err, domain.ErrInvalidCatalogKind)
}

func TestRepositoryListReturnsUnavailableWhenCatalogWasNotLoaded(t *testing.T) {
	repository := memory.NewRepository()

	_, err := repository.List(context.Background(), domain.CatalogDuty)

	require.ErrorIs(t, err, domain.ErrCatalogUnavailable)
}

func TestRepositoryGetByIDReturnsPharmacy(t *testing.T) {
	repository := memory.NewRepository()
	require.NoError(t, repository.Replace(context.Background(), domain.CatalogRegular, []domain.Pharmacy{
		{ID: "1", Name: "Cruz Verde"},
		{ID: "", Name: "Without ID"},
	}))

	pharmacy, err := repository.GetByID(context.Background(), domain.CatalogRegular, " 1 ")

	require.NoError(t, err)
	assert.Equal(t, "1", pharmacy.ID)
	assert.Equal(t, "Cruz Verde", pharmacy.Name)
}

func TestRepositoryGetByIDReturnsErrors(t *testing.T) {
	repository := memory.NewRepository()

	_, err := repository.GetByID(context.Background(), domain.CatalogRegular, "1")
	require.ErrorIs(t, err, domain.ErrCatalogUnavailable)

	require.NoError(t, repository.Replace(context.Background(), domain.CatalogRegular, []domain.Pharmacy{{ID: "1", Name: "Cruz Verde"}}))

	_, err = repository.GetByID(context.Background(), domain.CatalogRegular, "missing")
	require.ErrorIs(t, err, domain.ErrPharmacyNotFound)

	_, err = repository.GetByID(context.Background(), domain.CatalogRegular, "")
	require.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestRepositoryGetByIDRejectsInvalidCatalog(t *testing.T) {
	repository := memory.NewRepository()

	_, err := repository.GetByID(context.Background(), domain.CatalogKind("bad"), "1")

	require.ErrorIs(t, err, domain.ErrInvalidCatalogKind)
}

func TestRepositoryReplaceRejectsInvalidCatalog(t *testing.T) {
	repository := memory.NewRepository()

	err := repository.Replace(context.Background(), domain.CatalogKind("bad"), nil)

	require.ErrorIs(t, err, domain.ErrInvalidCatalogKind)
}
