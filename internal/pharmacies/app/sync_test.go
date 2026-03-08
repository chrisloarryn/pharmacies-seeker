package app_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
)

type providerStub struct {
	data map[domain.CatalogKind][]domain.Pharmacy
	errs map[domain.CatalogKind]error
}

func (s providerStub) Fetch(_ context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	if err := kind.Validate(); err != nil {
		return nil, err
	}
	if err := s.errs[kind]; err != nil {
		return nil, err
	}
	return append([]domain.Pharmacy(nil), s.data[kind]...), nil
}

func TestSyncServiceSyncAllLoadsRegularAndDutyCatalogs(t *testing.T) {
	repository := memory.NewRepository()
	service := app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "1", Name: "Regular"}},
			domain.CatalogDuty:    {{ID: "2", Name: "Duty"}},
		},
		errs: map[domain.CatalogKind]error{},
	}, repository)

	err := service.SyncAll(context.Background())

	require.NoError(t, err)
	assert.True(t, service.Ready())

	regular, err := repository.List(context.Background(), domain.CatalogRegular)
	require.NoError(t, err)
	require.Len(t, regular, 1)

	duty, err := repository.List(context.Background(), domain.CatalogDuty)
	require.NoError(t, err)
	require.Len(t, duty, 1)
}

func TestSyncServiceKeepsPreviousSnapshotWhenRefreshFails(t *testing.T) {
	repository := memory.NewRepository()
	service := app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "regular-v1", Name: "Regular v1"}},
			domain.CatalogDuty:    {{ID: "duty-v1", Name: "Duty v1"}},
		},
		errs: map[domain.CatalogKind]error{},
	}, repository)

	require.NoError(t, service.SyncAll(context.Background()))

	service = app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "regular-v2", Name: "Regular v2"}},
		},
		errs: map[domain.CatalogKind]error{
			domain.CatalogDuty: errors.New("boom"),
		},
	}, repository)
	service.SyncAll(context.Background())

	regular, err := repository.List(context.Background(), domain.CatalogRegular)
	require.NoError(t, err)
	assert.Equal(t, "regular-v2", regular[0].ID)

	duty, err := repository.List(context.Background(), domain.CatalogDuty)
	require.NoError(t, err)
	assert.Equal(t, "duty-v1", duty[0].ID)
}

type repositoryStub struct {
	replaceErrs map[domain.CatalogKind]error
}

func (s repositoryStub) Replace(_ context.Context, kind domain.CatalogKind, _ []domain.Pharmacy) error {
	return s.replaceErrs[kind]
}

func (s repositoryStub) List(context.Context, domain.CatalogKind) ([]domain.Pharmacy, error) {
	return nil, nil
}

func (s repositoryStub) GetByID(context.Context, domain.CatalogKind, string) (domain.Pharmacy, error) {
	return domain.Pharmacy{}, nil
}

func TestSyncServiceReturnsErrorAndStaysNotReadyWhenBothFetchesFail(t *testing.T) {
	service := app.NewSyncService(providerStub{
		errs: map[domain.CatalogKind]error{
			domain.CatalogRegular: errors.New("regular failed"),
			domain.CatalogDuty:    errors.New("duty failed"),
		},
	}, repositoryStub{})

	err := service.SyncAll(context.Background())

	require.Error(t, err)
	assert.False(t, service.Ready())
	assert.Contains(t, err.Error(), "regular catalog fetch")
	assert.Contains(t, err.Error(), "duty catalog fetch")
}

func TestSyncServiceReturnsReplaceError(t *testing.T) {
	service := app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "1"}},
			domain.CatalogDuty:    {{ID: "2"}},
		},
	}, repositoryStub{
		replaceErrs: map[domain.CatalogKind]error{
			domain.CatalogDuty: errors.New("cannot save"),
		},
	})

	err := service.SyncAll(context.Background())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "duty catalog replace")
	assert.False(t, service.Ready())
}
