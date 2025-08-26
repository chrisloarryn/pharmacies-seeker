package fetcherService

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
)

type apiClientHTTP struct {
	url     string
	dutyURL string
}

type boostrResponse struct {
	Status string           `json:"status"`
	Data   []boostrPharmacy `json:"data"`
}

type boostrPharmacy struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (c *apiClientHTTP) fetch(ctx context.Context, url string) ([]pharmacy.Pharmacy, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("remote API returned status %d: %s", resp.StatusCode, string(bytes.TrimSpace(b)))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// The error referring to ï is because the UTF-8 BOM interpreted as an ISO-8859-1 string will produce the characters ï»¿.
	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf"))

	var br boostrResponse
	if err := json.Unmarshal(bodyBytes, &br); err != nil {
		return nil, err
	}

	// Map Boostr data to domain Pharmacy
	out := make([]pharmacy.Pharmacy, 0, len(br.Data))
	for _, p := range br.Data {
		out = append(out, pharmacy.Pharmacy{
			LocalNombre:    p.Name,
			ComunaNombre:   p.City,
			LocalDireccion: p.Street,
			LocalTelefono:  p.Phone,
		})
	}
	return out, nil
}

func (c *apiClientHTTP) RetrievePharmacies(ctx context.Context) ([]pharmacy.Pharmacy, error) {
	return c.fetch(ctx, c.url)
}

func (c *apiClientHTTP) RetrievePharmacies24h(ctx context.Context) ([]pharmacy.Pharmacy, error) {
	return c.fetch(ctx, c.dutyURL)
}

func NewFetcherService(c config.Config) fetcher.Service {
	return &apiClientHTTP{url: c.Api.Pharmacy.Url, dutyURL: c.Api.Pharmacy.DutyUrl}
}
