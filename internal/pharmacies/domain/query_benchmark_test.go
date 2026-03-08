package domain_test

import (
	"testing"

	"pharmacies-seeker/internal/pharmacies/domain"
)

func BenchmarkMatches(b *testing.B) {
	pharmacy := domain.Pharmacy{
		ID:      "1",
		Name:    "Cruz Verde Alameda",
		Commune: "Santiago Centro",
		Address: "Alameda 100",
		Phone:   "123",
	}
	query := domain.NewQuery("santiago", "verde")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = domain.Matches(pharmacy, query)
	}
}
