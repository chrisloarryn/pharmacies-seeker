package pharmacyhttp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pharmacyhttp "pharmacies-seeker/internal/pharmacies/adapters/http"
	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
)

type providerStub struct {
	data map[domain.CatalogKind][]domain.Pharmacy
}

func (s providerStub) Fetch(_ context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	return append([]domain.Pharmacy(nil), s.data[kind]...), nil
}

func TestListPharmaciesReturnsOK(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies?commune=santiago&name=verde", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestLiveEndpointReturnsOK(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestReadyEndpointReturnsOKAfterStartupSync(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetPharmacyByIDReturnsNotFound(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/missing", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetPharmacyByIDReturnsOK(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/1", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestGetPharmacyByIDReturnsBadRequestForBlankID(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/%E2%98%83", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestReadyEndpointReflectsSyncState(t *testing.T) {
	repository := memory.NewRepository()
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(providerStub{data: map[domain.CatalogKind][]domain.Pharmacy{}}, repository)
	router := pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestCatalogEndpointsReturnServiceUnavailableBeforeStartupSync(t *testing.T) {
	repository := memory.NewRepository()
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(providerStub{data: map[domain.CatalogKind][]domain.Pharmacy{}}, repository)
	router := pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestGetPharmacyByIDReturnsServiceUnavailableBeforeStartupSync(t *testing.T) {
	repository := memory.NewRepository()
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(providerStub{data: map[domain.CatalogKind][]domain.Pharmacy{}}, repository)
	router := pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/1", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestListDutyPharmaciesReturnsOK(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/duty?commune=provi", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestListDutyPharmaciesReturnsServiceUnavailableBeforeStartupSync(t *testing.T) {
	repository := memory.NewRepository()
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(providerStub{data: map[domain.CatalogKind][]domain.Pharmacy{}}, repository)
	router := pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/duty", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

type failingRepositoryStub struct {
	listErr error
	getErr  error
}

func (s failingRepositoryStub) Replace(context.Context, domain.CatalogKind, []domain.Pharmacy) error {
	return nil
}

func (s failingRepositoryStub) List(context.Context, domain.CatalogKind) ([]domain.Pharmacy, error) {
	return nil, s.listErr
}

func (s failingRepositoryStub) GetByID(context.Context, domain.CatalogKind, string) (domain.Pharmacy, error) {
	return domain.Pharmacy{}, s.getErr
}

func TestListPharmaciesReturnsInternalErrorForUnexpectedRepositoryFailure(t *testing.T) {
	router := newRouterWithCatalogs(t, app.NewCatalogService(failingRepositoryStub{listErr: errors.New("boom")}))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

func TestGetPharmacyByIDReturnsServiceUnavailableWhenCatalogIsMissingAfterReadySync(t *testing.T) {
	router := newRouterWithCatalogs(t, app.NewCatalogService(memory.NewRepository()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/1", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func TestListDutyPharmaciesReturnsServiceUnavailableWhenCatalogIsMissingAfterReadySync(t *testing.T) {
	router := newRouterWithCatalogs(t, app.NewCatalogService(memory.NewRepository()))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/duty", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)
}

func newReadyRouter(t *testing.T) *fiber.App {
	t.Helper()

	repository := memory.NewRepository()
	catalogs := app.NewCatalogService(repository)
	syncService := app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {
				{ID: "1", Name: "Cruz Verde", Commune: "Santiago", Address: "Street 1", Phone: "123"},
			},
			domain.CatalogDuty: {
				{ID: "2", Name: "Salcobrand", Commune: "Providencia", Address: "Street 2", Phone: "456"},
			},
		},
	}, repository)

	require.NoError(t, syncService.SyncAll(context.Background()))

	return pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))
}

func newRouterWithCatalogs(t *testing.T, catalogs *app.CatalogService) *fiber.App {
	t.Helper()

	repository := memory.NewRepository()
	syncService := app.NewSyncService(providerStub{
		data: map[domain.CatalogKind][]domain.Pharmacy{
			domain.CatalogRegular: {{ID: "ready"}},
			domain.CatalogDuty:    {{ID: "ready-duty"}},
		},
	}, repository)

	require.NoError(t, syncService.SyncAll(context.Background()))

	return pharmacyhttp.NewRouter(pharmacyhttp.NewHandler(catalogs, syncService))
}
