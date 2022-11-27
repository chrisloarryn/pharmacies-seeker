package dependencies

import (
	"pharmacies-seeker/cmd/config"
	"pharmacies-seeker/internal/core/domain/fetcher"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/http/client/fetcherService"
	"pharmacies-seeker/internal/infraestucture/storage"
)

type Container interface {
	Config() config.Config
	FetcherService() fetcher.Service
	PharmaciesRepository() pharmacy.Repository
}

type container struct {
	cfg                  config.Config
	fetcherService       fetcher.Service
	pharmaciesRepository pharmacy.Repository
}

func NewContainer(c config.Config) Container {
	fSvc := fetcherService.NewFetcherService(c)

	return &container{
		cfg:                  c,
		pharmaciesRepository: storage.New(c),
		fetcherService:       fSvc,
	}
}

func (c *container) Config() config.Config {
	return c.cfg
}

func (c *container) FetcherService() fetcher.Service {
	return c.fetcherService
}

func (c *container) PharmaciesRepository() pharmacy.Repository {
	return c.pharmaciesRepository
}
