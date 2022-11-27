package dependencies

import (
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/http/client/fetcherService"
	"pharmacies-seeker/internal/infraestucture/storage"
)

type Container interface {
	FetcherService() fetcher.Service
	PharmaciesRepository() pharmacy.Repository
}

type container struct {
	fetcherService      fetcher.Service
	pharmaciesRepository pharmacy.Repository
}

func NewContainer() Container {
	fSvc := fetcherService.NewFetcherService()
	return &container{
		pharmaciesRepository: storage.New(),
		fetcherService:      fSvc,
	}
}

func (c *container) FetcherService() fetcher.Service {
	return c.fetcherService
}

func (c *container) PharmaciesRepository() pharmacy.Repository {
	return c.pharmaciesRepository
}
