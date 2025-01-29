package geocoding

import (
	"errors"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
)

var (
	ErrGeocodeProvidersUnavailable = errors.New("geocode providers unavailable")
)

type GeocodeConfig struct {
	Lang  language.Tag
	Limit int
}

type ReverseGeocodeConfig struct {
	Lang  language.Tag
	Limit int
	Zoom  int
}

type Location struct {
	Lat decimal.Decimal
	Lng decimal.Decimal
}

type Address struct {
	FormattedAddress string
	Street           string
	HouseNumber      string
	Suburb           string
	Postcode         string
	State            string
	StateCode        string
	StateDistrict    string
	County           string
	Country          string
	CountryCode      string
	City             string
}
