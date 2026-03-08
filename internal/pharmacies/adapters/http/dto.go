package pharmacyhttp

import "pharmacies-seeker/internal/pharmacies/domain"

type errorEnvelope struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type pharmacyDTO struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Commune string `json:"commune"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type pharmacyListResponse struct {
	Data []pharmacyDTO `json:"data"`
}

type pharmacyResponse struct {
	Data pharmacyDTO `json:"data"`
}

func toPharmacyDTO(pharmacy domain.Pharmacy) pharmacyDTO {
	return pharmacyDTO{
		ID:      pharmacy.ID,
		Name:    pharmacy.Name,
		Commune: pharmacy.Commune,
		Address: pharmacy.Address,
		Phone:   pharmacy.Phone,
	}
}

func newPharmacyListResponse(pharmacies []domain.Pharmacy) pharmacyListResponse {
	items := make([]pharmacyDTO, 0, len(pharmacies))
	for _, pharmacy := range pharmacies {
		items = append(items, toPharmacyDTO(pharmacy))
	}

	return pharmacyListResponse{Data: items}
}

func newPharmacyResponse(pharmacy domain.Pharmacy) pharmacyResponse {
	return pharmacyResponse{Data: toPharmacyDTO(pharmacy)}
}
