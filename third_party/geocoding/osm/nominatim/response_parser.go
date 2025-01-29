package nominatim

import (
	"errors"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/shopspring/decimal"
	"strings"
)

type geocodeResponse []struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type reverseGeocodeResponse struct {
	DisplayName string  `json:"display_name"`
	Addr        Address `json:"address"`
}

type errorOutput struct {
	ErrorStr string `json:"error"`
}

type responseParser struct {
	geocodeResponse
	reverseGeocodeResponse
	errorOutput
}

func (r *geocodeResponse) Locations() ([]*geocoding.Location, error) {
	locs := make([]*geocoding.Location, len(*r))
	for i, resp := range *r {
		if resp.Lat == "" || resp.Lon == "" {
			continue
		}

		lng, err := decimal.NewFromString(resp.Lon)
		if err != nil {
			continue
		}

		lat, err := decimal.NewFromString(resp.Lat)
		if err != nil {
			continue
		}

		locs[i] = &geocoding.Location{
			Lat: lat,
			Lng: lng,
		}
	}

	return locs, nil
}

func (r *reverseGeocodeResponse) Addresses() ([]*geocoding.Address, error) {
	addr := &geocoding.Address{
		FormattedAddress: r.DisplayName,
		Street:           r.Addr.Street(),
		HouseNumber:      r.Addr.HouseNumber,
		City:             r.Addr.Locality(),
		Postcode:         r.Addr.Postcode,
		Suburb:           r.Addr.Suburb,
		State:            r.Addr.State,
		Country:          r.Addr.Country,
		CountryCode:      strings.ToUpper(r.Addr.CountryCode),
	}

	return []*geocoding.Address{addr}, nil
}

func (r *errorOutput) Error() error {
	return errors.New(r.ErrorStr)
}
