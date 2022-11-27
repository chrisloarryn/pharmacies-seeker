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

func TestFinderAllPharmacies_Execute_ShouldReturnsAPharmacyList(t *testing.T) {
	t.Log("Should returns a pharmacy list")
	// Setup
	controller := gomock.NewController(t)

	repository := mocks.NewMockRepository(controller)

	communeName := "test commune name"
	pharmaciesList := []pharmacy.Pharmacy{{
		LocalNombre:    "test local name",
		ComunaNombre:   "test commune name",
		LocalDireccion: "test address",
		LocalTelefono:  "test phone",
	}, {
		LocalNombre:    "test local name 2",
		ComunaNombre:   "test commune name 2",
		LocalDireccion: "test address 2",
		LocalTelefono:  "test phone 2",
	}}
	repository.EXPECT().Find(gomock.Any(), communeName).Return(pharmaciesList, nil).Times(1)

	find := usecases.NewFinderAllPharmacies(repository)

	// Execute
	result, err := find.Execute(context.TODO(), communeName)

	// Verify
	require.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, pharmaciesList, result)
}
