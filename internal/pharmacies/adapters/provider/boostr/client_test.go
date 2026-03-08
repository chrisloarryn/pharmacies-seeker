package boostr_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/adapters/provider/boostr"
	"pharmacies-seeker/internal/pharmacies/domain"
)

func TestClientFetchMapsRemotePayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/regular", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "\ufeff{\"status\":\"ok\",\"data\":[{\"id\":\"1\",\"name\":\"Cruz Verde\",\"phone\":\"123\",\"street\":\"Main\",\"city\":\"Santiago\"}]}")
	}))
	defer server.Close()

	client := boostr.NewClient(server.URL+"/regular", server.URL+"/duty", time.Second)

	pharmacies, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.NoError(t, err)
	require.Len(t, pharmacies, 1)
	assert.Equal(t, "1", pharmacies[0].ID)
	assert.Equal(t, "Cruz Verde", pharmacies[0].Name)
}

func TestClientFetchReturnsErrorForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "upstream failed", http.StatusBadGateway)
	}))
	defer server.Close()

	client := boostr.NewClient(server.URL, server.URL, time.Second)

	_, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "status 502")
}

func TestClientFetchReturnsErrorForInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{invalid")
	}))
	defer server.Close()

	client := boostr.NewClient(server.URL, server.URL, time.Second)

	_, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.Error(t, err)
}

func TestClientFetchUsesDutyURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/duty", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, "{\"status\":\"ok\",\"data\":[]}")
	}))
	defer server.Close()

	client := boostr.NewClient(server.URL+"/regular", server.URL+"/duty", time.Second)

	pharmacies, err := client.Fetch(context.Background(), domain.CatalogDuty)

	require.NoError(t, err)
	assert.Empty(t, pharmacies)
}

func TestClientFetchRejectsInvalidCatalogKind(t *testing.T) {
	client := boostr.NewClient("http://example.com", "http://example.com", time.Second)

	_, err := client.Fetch(context.Background(), domain.CatalogKind("bad"))

	require.ErrorIs(t, err, domain.ErrInvalidCatalogKind)
}

func TestClientFetchReturnsErrorWhenRequestCannotBeCreated(t *testing.T) {
	client := boostr.NewClient("://bad-url", "http://example.com", time.Second)

	_, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.Error(t, err)
}

func TestClientFetchReturnsTransportError(t *testing.T) {
	client := boostr.NewClient("http://127.0.0.1:1", "http://127.0.0.1:1", 50*time.Millisecond)

	_, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.Error(t, err)
}
