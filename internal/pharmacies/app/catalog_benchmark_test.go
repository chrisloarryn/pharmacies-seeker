package app_test

import (
	"context"
	"fmt"
	"testing"

	"pharmacies-seeker/internal/pharmacies/adapters/repository/memory"
	"pharmacies-seeker/internal/pharmacies/app"
	"pharmacies-seeker/internal/pharmacies/domain"
)

func BenchmarkCatalogServiceList(b *testing.B) {
	repository := memory.NewRepository()
	items := make([]domain.Pharmacy, 0, 1000)
	for i := 0; i < 1000; i++ {
		items = append(items, domain.Pharmacy{
			ID:      fmt.Sprintf("%d", i),
			Name:    fmt.Sprintf("Pharmacy %d", i),
			Commune: "Santiago",
			Address: "Address",
			Phone:   "123",
		})
	}
	_ = repository.Replace(context.Background(), domain.CatalogRegular, items)

	service := app.NewCatalogService(repository)
	query := domain.NewQuery("santiago", "Pharmacy")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.List(context.Background(), domain.CatalogRegular, query)
	}
}
