package domain

import "errors"

var (
	ErrInvalidCatalogKind = errors.New("invalid catalog kind")
	ErrInvalidID          = errors.New("invalid pharmacy id")
	ErrCatalogUnavailable = errors.New("catalog not synchronized yet")
	ErrPharmacyNotFound   = errors.New("pharmacy not found")
)
