package fetcherService

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/shared/constants"
)

type apiClientHTTP struct {
	url string
}

func (c *apiClientHTTP) RetrievePharmacies(ctx context.Context) ([]pharmacy.Pharmacy, error) {
	url := c.url

	result, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	// The error referring to ï is because the UTF-8 BOM interpreted as an ISO-8859-1 string will produce the characters ï»¿.
	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}

	var pharmacies []pharmacy.Pharmacy

	err = json.Unmarshal(bodyBytes, &pharmacies)
	if err != nil {
		return nil, err
	}

	return pharmacies, nil
}

func NewFetcherService() fetcher.Service {
	serviceURL := os.Getenv(constants.ApiServiceUrlName)
	return &apiClientHTTP{url: serviceURL}
}
