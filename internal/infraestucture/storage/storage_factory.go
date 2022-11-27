package storage

import (
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/infraestucture/storage/inmemory"
)

func New(c config.Config) pharmacy.Repository {
	return inmemory.NewInMemoryRepository(c)
}
