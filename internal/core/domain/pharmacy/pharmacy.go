package pharmacy

import (
	"context"
	"encoding/json"
	"errors"
)

func UnmarshalPharmacies(data []byte) (Pharmacies, error) {
	var r Pharmacies
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Pharmacies) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Pharmacy) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Pharmacy) UnmarshalToXML(result interface{}) error {
	b, err := r.Marshal()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, result)
}

func (r *Pharmacy) ToXMLInterface() (PharmacyXML, error) {
	var xml PharmacyXML
	err := r.UnmarshalToXML(&xml)
	return xml, err
}

// Pharmacies is an array of Pharmacy
//
// This type is used for unmarshaling JSON.
//
// swagger:model Pharmacy
type Pharmacies struct {
	// The pharmacies
	// in: body
	Pharmacies []Pharmacy `json:"pharmacies,omitempty" xml:"PHARMACIES,omitempty"`
}

// Pharmacy is a pharmacy
//
// This type is used for unmarshaling JSON.
//
// swagger:model Pharmacy
type Pharmacy struct {
	// The local name
	LocalNombre string `json:"local_nombre,omitempty" xml:"LOCAL_NOMBRE,omitempty"`
	// The commune name
	ComunaNombre string `json:"comuna_nombre,omitempty" xml:"COMUNA_NOMBRE,omitempty"`
	// The local address
	LocalDireccion string `json:"local_direccion,omitempty" xml:"LOCAL_DIRECCION,omitempty"`
	// The local phone
	LocalTelefono string `json:"local_telefono,omitempty" xml:"LOCAL_TELEFONO,omitempty"`
}

// PharmacyXML is a pharmacy
//
// This type is used for unmarshaling JSON.
//
// swagger:model PharmacyXML
type PharmacyXML struct {
	// The local name
	LocalNombre string `xml:"LOCAL_NOMBRE,omitempty"`
	// The commune name
	ComunaNombre string `xml:"COMUNA_NOMBRE,omitempty"`
	// The local address
	LocalDireccion string `xml:"LOCAL_DIRECCION,omitempty"`
	// The local phone
	LocalTelefono string `xml:"LOCAL_TELEFONO,omitempty"`
}

//go:generate mockgen -package mocks -destination mocks/pharmacies_repository_mocks.go . Repository

func ValidatePharmacyCommune(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}
	return nil
}

// Repository is the storage abstraction
type Repository interface {
	Find(ctx context.Context, commune string) ([]Pharmacy, error)
	FindOne(ctx context.Context, commune string) (Pharmacy, error)
	LoadAll(ctx context.Context, pharmacies []Pharmacy) error
}
