package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCatalogKindValidateAcceptsKnownValues(t *testing.T) {
	require.NoError(t, CatalogRegular.Validate())
	require.NoError(t, CatalogDuty.Validate())
}

func TestCatalogKindValidateRejectsUnknownValue(t *testing.T) {
	err := CatalogKind("other").Validate()

	require.ErrorIs(t, err, ErrInvalidCatalogKind)
}
