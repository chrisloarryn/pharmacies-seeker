package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pharmacies-seeker/internal/pharmacies/domain"
)

func TestNewQueryTrimsInput(t *testing.T) {
	query := domain.NewQuery("  Santiago ", " Cruz Verde  ")

	assert.Equal(t, "Santiago", query.Commune)
	assert.Equal(t, "Cruz Verde", query.Name)
	require.NoError(t, query.Validate())
}

func TestNormalizeIDRejectsBlankValue(t *testing.T) {
	_, err := domain.NormalizeID("   ")

	require.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestNormalizeIDRejectsUnsupportedCharacters(t *testing.T) {
	_, err := domain.NormalizeID("snowman-☃")

	require.ErrorIs(t, err, domain.ErrInvalidID)
}

func TestNormalizeIDAcceptsSupportedCharacters(t *testing.T) {
	id, err := domain.NormalizeID(" Abc-123_. ")

	require.NoError(t, err)
	assert.Equal(t, "Abc-123_.", id)
}

func TestMatchesUsesCaseInsensitiveContains(t *testing.T) {
	pharmacy := domain.Pharmacy{
		ID:      "1",
		Name:    "Cruz Verde Central",
		Commune: "Santiago Centro",
	}

	assert.True(t, domain.Matches(pharmacy, domain.NewQuery("santiago", "verde")))
	assert.False(t, domain.Matches(pharmacy, domain.NewQuery("valparaiso", "")))
}

func TestMatchesRejectsNameWhenItDoesNotMatch(t *testing.T) {
	pharmacy := domain.Pharmacy{
		ID:      "1",
		Name:    "Cruz Verde Central",
		Commune: "Santiago Centro",
	}

	assert.False(t, domain.Matches(pharmacy, domain.NewQuery("santiago", "ahumada")))
}
