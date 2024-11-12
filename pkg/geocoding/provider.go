package geocoding

import (
	"context"
	"golang.org/x/text/language"
)

type Location struct {
	Lat float64
	Lng float64
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

type GeocodeConfig struct {
	Lang  language.Tag
	Limit int
}

type ReverseGeocodeConfig struct {
	Lang  language.Tag
	Limit int
	Zoom  int
}

type ProviderError struct {
	status int
	msg    string
}

func NewProviderError(msg string, status int) *ProviderError {
	return &ProviderError{msg: msg, status: status}
}

func (e *ProviderError) Error() string {
	return e.msg
}

func (e *ProviderError) Status() int {
	return e.status
}

type Provider interface {
	Geocode(address string, config GeocodeConfig) ([]*Location, error)
	GeocodeWithContext(ctx context.Context, address string, config GeocodeConfig) ([]*Location, error)
	ReverseGeocode(loc Location, config ReverseGeocodeConfig) ([]*Address, error)
	ReverseGeocodeWithContext(ctx context.Context, loc Location, config ReverseGeocodeConfig) ([]*Address, error)
}
