package domain

type Pharmacy struct {
	ID      string
	Name    string
	Commune string
	Address string
	Phone   string
}

type CatalogKind string

const (
	CatalogRegular CatalogKind = "regular"
	CatalogDuty    CatalogKind = "duty"
)

func (k CatalogKind) Validate() error {
	switch k {
	case CatalogRegular, CatalogDuty:
		return nil
	default:
		return ErrInvalidCatalogKind
	}
}
