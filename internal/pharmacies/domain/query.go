package domain

import "strings"

type Query struct {
	Commune string
	Name    string
}

func NewQuery(commune, name string) Query {
	return Query{
		Commune: strings.TrimSpace(commune),
		Name:    strings.TrimSpace(name),
	}
}

func (q Query) Validate() error {
	return nil
}

func NormalizeID(id string) (string, error) {
	normalized := strings.TrimSpace(id)
	if normalized == "" {
		return "", ErrInvalidID
	}
	for _, char := range normalized {
		switch {
		case char >= 'a' && char <= 'z':
		case char >= 'A' && char <= 'Z':
		case char >= '0' && char <= '9':
		case char == '-', char == '_', char == '.':
		default:
			return "", ErrInvalidID
		}
	}
	return normalized, nil
}

func Matches(pharmacy Pharmacy, query Query) bool {
	query = NewQuery(query.Commune, query.Name)

	if query.Commune != "" && !containsFold(pharmacy.Commune, query.Commune) {
		return false
	}
	if query.Name != "" && !containsFold(pharmacy.Name, query.Name) {
		return false
	}
	return true
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
