package app

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"pharmacies-seeker/internal/pharmacies/domain"
)

type SyncService struct {
	provider   Provider
	repository Repository
	ready      atomic.Bool
}

func NewSyncService(provider Provider, repository Repository) *SyncService {
	return &SyncService{
		provider:   provider,
		repository: repository,
	}
}

func (s *SyncService) SyncAll(ctx context.Context) error {
	var syncErr error

	if err := s.syncCatalog(ctx, domain.CatalogRegular); err != nil {
		syncErr = err
	}
	if err := s.syncCatalog(ctx, domain.CatalogDuty); err != nil {
		if syncErr == nil {
			syncErr = err
		} else {
			syncErr = errors.Join(syncErr, err)
		}
	}

	if syncErr == nil {
		s.ready.Store(true)
	}

	return syncErr
}

func (s *SyncService) Ready() bool {
	return s.ready.Load()
}

func (s *SyncService) syncCatalog(ctx context.Context, kind domain.CatalogKind) error {
	pharmacies, err := s.provider.Fetch(ctx, kind)
	if err != nil {
		return fmt.Errorf("%s catalog fetch: %w", kind, err)
	}
	if err := s.repository.Replace(ctx, kind, pharmacies); err != nil {
		return fmt.Errorf("%s catalog replace: %w", kind, err)
	}
	return nil
}
