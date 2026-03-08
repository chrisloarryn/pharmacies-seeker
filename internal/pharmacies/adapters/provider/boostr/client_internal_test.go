package boostr

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type brokenBody struct{}

func (brokenBody) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (brokenBody) Close() error {
	return nil
}

func TestClientFetchReturnsErrorWhenResponseBodyCannotBeRead(t *testing.T) {
	client := NewClient("https://example.com/regular", "https://example.com/duty", time.Second)
	client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       brokenBody{},
			Header:     make(http.Header),
			Request:    req,
		}, nil
	})

	_, err := client.Fetch(context.Background(), domain.CatalogRegular)

	require.EqualError(t, err, "read failed")
}
