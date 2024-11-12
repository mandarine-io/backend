package nominatim

import (
	"errors"
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"strconv"
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

type errorResponse struct {
	ErrorStr string `json:"error"`
}

type responseParser struct {
	geocodeResponse
	reverseGeocodeResponse
	errorResponse
}

func (r *geocodeResponse) Locations() ([]*geocoding.Location, error) {
	locs := make([]*geocoding.Location, len(*r))
	for i, resp := range *r {
		if resp.Lat == "" || resp.Lon == "" {
			continue
		}

		lat, err := strconv.ParseFloat(resp.Lat, 64)
		if err != nil {
			continue
		}
		lon, err := strconv.ParseFloat(resp.Lon, 64)
		if err != nil {
			continue
		}

		locs[i] = &geocoding.Location{
			Lat: lat,
			Lng: lon,
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

func (r *errorResponse) Error() error {
	return errors.New(r.ErrorStr)
}
