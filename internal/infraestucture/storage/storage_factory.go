package storage

import (
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/infraestucture/storage/inmemory"
)

func New() pharmacy.Repository {
	return inmemory.NewInMemoryRepository()
}
