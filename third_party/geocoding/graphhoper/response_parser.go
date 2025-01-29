package graphhopper

import (
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/shopspring/decimal"
	"reflect"
	"strings"
)

type geocodeHint struct {
	Point struct {
		Lat decimal.Decimal `json:"lat"`
		Lng decimal.Decimal `json:"lng"`
	} `json:"point"`
}

type reverseGeocodeHint struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	State       string `json:"state"`
	Street      string `json:"street"`
	HouseNumber string `json:"housenumber"`
	PostCode    string `json:"postcode"`
}

type geocodeResponse struct {
	Hits []geocodeHint `json:"hits"`
}

type reverseGeocodeResponse struct {
	Hits []reverseGeocodeHint `json:"hits"`
}

type responseParser struct {
	geocodeResponse
	reverseGeocodeResponse
}

func (r *geocodeResponse) Locations() ([]*geocoding.Location, error) {
	if len(r.Hits) == 0 {
		return []*geocoding.Location{}, nil
	}

	result := make([]*geocoding.Location, len(r.Hits))
	for i, hit := range r.Hits {
		result[i] = &geocoding.Location{
			Lat: hit.Point.Lat,
			Lng: hit.Point.Lng,
		}
	}

	return result, nil
}

func (r *reverseGeocodeResponse) Addresses() ([]*geocoding.Address, error) {
	if len(r.Hits) == 0 {
		return []*geocoding.Address{}, nil
	}

	addrs := make([]*geocoding.Address, len(r.Hits))
	for i, hit := range r.Hits {
		addrs[i] = &geocoding.Address{
			FormattedAddress: formatAddress(hit),
			Street:           hit.Street,
			HouseNumber:      hit.HouseNumber,
			City:             hit.City,
			Postcode:         hit.PostCode,
			State:            hit.State,
			Country:          hit.Country,
			CountryCode:      "",
			Suburb:           "",
			StateCode:        "",
			StateDistrict:    "",
			County:           "",
		}
	}

	return addrs, nil
}

func formatAddress(hint reverseGeocodeHint) string {
	parts := []string{
		hint.Country,
		hint.PostCode,
		hint.State,
		hint.City,
		hint.Street,
		hint.HouseNumber,
	}

	filtered := make([]string, 0)
	for _, part := range parts {
		if !reflect.ValueOf(part).IsZero() {
			filtered = append(filtered, part)
		}
	}

	return strings.Join(filtered, ", ")
}

func (r *responseParser) Error() error {
	return nil
}
