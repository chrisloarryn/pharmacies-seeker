package usecases_test

import (
	"context"
	"pharmacies-seeker/internal/core/domain/pharmacy"
	"pharmacies-seeker/internal/core/domain/pharmacy/mocks"
	"pharmacies-seeker/internal/core/usecases"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinderPharmaciesByCommune_Execute_ShouldReturnsAPharmacyData(t *testing.T) {
	t.Log("Should returns a pharmacy from his commune name")
	// Setup
	controller := gomock.NewController(t)

	communeName := "test commune name"
	pharmacyResult := pharmacy.Pharmacy{
		LocalNombre:    "test local name",
		ComunaNombre:   "test commune name",
		LocalDireccion: "test address",
		LocalTelefono:  "test phone",
	}

	repository := mocks.NewMockRepository(controller)
	repository.EXPECT().FindOne(gomock.Any(), communeName).Return(pharmacyResult, nil).Times(1)

	find := usecases.NewFinderPharmaciesByCommune(repository)

	// Execute
	result, err := find.Execute(context.TODO(), communeName)

	// Verify
	require.NoError(t, err)
	assert.Equal(t, pharmacyResult, result)
}

func TestFinderPharmacyByCommuneName_Execute_ShouldReturnsAnErrorForInvalidName(t *testing.T) {
	t.Log("Should returns an error for invalid commune name")
	// Setup
	controller := gomock.NewController(t)

	communeName := ""
	repository := mocks.NewMockRepository(controller)

	find := usecases.NewFinderPharmaciesByCommune(repository)

	// Execute
	result, err := find.Execute(context.TODO(), communeName)

	// Verify
	require.Error(t, err, "Invalid name")
	assert.Equal(t, pharmacy.Pharmacy{}, result)
}
