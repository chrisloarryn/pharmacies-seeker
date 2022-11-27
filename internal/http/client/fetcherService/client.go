package fetcherService

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
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

	bodyBytes = bytes.TrimPrefix(bodyBytes, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191}

	var pharmacies []pharmacy.Pharmacy

	err = json.Unmarshal(bodyBytes, &pharmacies)
	if err != nil {
		fmt.Println("ERROR:", url)
		return nil, err
	}

	return pharmacies, nil
}

func NewFetcherService(c config.Config) fetcher.Service {
	return &apiClientHTTP{url: c.Api.Pharmacy.Url}
}
