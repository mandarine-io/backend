package here

import (
	"errors"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/shopspring/decimal"
)

type geocodeResponse struct {
	Items []struct {
		Position struct {
			Lat decimal.Decimal
			Lng decimal.Decimal
		}
	}
}

type responseParserGeocodeResponse struct {
	Items []struct {
		Address struct {
			Label       string
			CountryCode string
			CountryName string
			StateCode   string
			State       string
			County      string
			District    string
			City        string
			Street      string
			PostalCode  string
			HouseNumber string
		}
	}
}

type errorOutput struct {
	Status           int    `json:"status"`
	Title            string `json:"title"`
	ErrorDescription string `json:"error_description"`
}

type responseParser struct {
	geocodeResponse
	responseParserGeocodeResponse
	errorOutput
}

func (r *geocodeResponse) Locations() ([]*geocoding.Location, error) {
	results := make([]*geocoding.Location, len(r.Items))
	for i, item := range r.Items {
		results[i] = &geocoding.Location{
			Lat: item.Position.Lat,
			Lng: item.Position.Lng,
		}
	}

	return results, nil
}

func (r *responseParserGeocodeResponse) Addresses() ([]*geocoding.Address, error) {
	results := make([]*geocoding.Address, len(r.Items))
	for i, item := range r.Items {
		res := item.Address
		results[i] = &geocoding.Address{
			FormattedAddress: res.Label,
			City:             res.City,
			Street:           res.Street,
			HouseNumber:      res.HouseNumber,
			Postcode:         res.PostalCode,
			State:            res.State,
			County:           res.County,
			Country:          res.CountryName,
			CountryCode:      res.CountryCode,
			Suburb:           res.District,
		}
	}
	return results, nil
}

func (r *errorOutput) Error() error {
	if r.Title != "" {
		return errors.New(r.Title)
	}
	return errors.New(r.ErrorDescription)
}
