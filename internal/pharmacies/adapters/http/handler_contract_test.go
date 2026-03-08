package pharmacyhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListPharmaciesContract(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies?commune=santiago&name=verde", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var body struct {
		Data []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Commune string `json:"commune"`
			Address string `json:"address"`
			Phone   string `json:"phone"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	require.Len(t, body.Data, 1)
	assert.Equal(t, "1", body.Data[0].ID)
	assert.Equal(t, "Cruz Verde", body.Data[0].Name)
}

func TestGetPharmacyByIDContract(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/1", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var body struct {
		Data struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Commune string `json:"commune"`
			Address string `json:"address"`
			Phone   string `json:"phone"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	assert.Equal(t, "1", body.Data.ID)
	assert.Equal(t, "Cruz Verde", body.Data.Name)
}

func TestListDutyPharmaciesContract(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/duty?commune=providencia", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var body struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	require.Len(t, body.Data, 1)
	assert.Equal(t, "2", body.Data[0].ID)
}

func TestErrorContract(t *testing.T) {
	router := newReadyRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/pharmacies/missing", nil)
	res, err := router.Test(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	assert.Equal(t, "not_found", body.Error.Code)
	assert.Equal(t, "pharmacy not found", body.Error.Message)
}
