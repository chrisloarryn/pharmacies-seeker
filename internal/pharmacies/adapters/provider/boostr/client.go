package boostr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type Client struct {
	regularURL string
	dutyURL    string
	httpClient *http.Client
}

type responsePayload struct {
	Status string            `json:"status"`
	Data   []pharmacyPayload `json:"data"`
}

type pharmacyPayload struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Street string `json:"street"`
	City   string `json:"city"`
}

func NewClient(regularURL, dutyURL string, timeout time.Duration) *Client {
	return &Client{
		regularURL: regularURL,
		dutyURL:    dutyURL,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) Fetch(ctx context.Context, kind domain.CatalogKind) ([]domain.Pharmacy, error) {
	if err := kind.Validate(); err != nil {
		return nil, err
	}

	url := c.regularURL
	if kind == domain.CatalogDuty {
		url = c.dutyURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("remote API returned status %d: %s", resp.StatusCode, string(bytes.TrimSpace(body)))
	}

	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))

	var payload responsePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	pharmacies := make([]domain.Pharmacy, 0, len(payload.Data))
	for _, pharmacy := range payload.Data {
		pharmacies = append(pharmacies, domain.Pharmacy{
			ID:      strings.TrimSpace(pharmacy.ID),
			Name:    strings.TrimSpace(pharmacy.Name),
			Commune: strings.TrimSpace(pharmacy.City),
			Address: strings.TrimSpace(pharmacy.Street),
			Phone:   strings.TrimSpace(pharmacy.Phone),
		})
	}

	return pharmacies, nil
}
