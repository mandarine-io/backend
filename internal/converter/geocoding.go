package converter

import (
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/shopspring/decimal"
)

func MapLngLatToLocation(lng, lat decimal.Decimal) geocoding.Location {
	return geocoding.Location{Lat: lat, Lng: lng}
}

func MapAddressToAddressOutput(address *geocoding.Address) v0.AddressOutput {
	return v0.AddressOutput{
		Country:          address.Country,
		CountryCode:      address.CountryCode,
		County:           address.County,
		City:             address.City,
		Postcode:         address.Postcode,
		State:            address.State,
		StateCode:        address.StateCode,
		StateDistrict:    address.StateDistrict,
		Street:           address.Street,
		HouseNumber:      address.HouseNumber,
		Suburb:           address.Suburb,
		FormattedAddress: address.FormattedAddress,
	}
}

func MapAddressesToGeocodingOutput(addresses []*geocoding.Address) v0.ReverseGeocodingOutput {
	data := make([]v0.AddressOutput, len(addresses))
	for i, address := range addresses {
		data[i] = MapAddressToAddressOutput(address)
	}

	return v0.ReverseGeocodingOutput{
		Count: len(addresses),
		Data:  data,
	}
}

func MapLocationToPointOutput(location *geocoding.Location) v0.PointOutput {
	return v0.PointOutput{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}
}

func MapLocationsToGeocodingOutput(locations []*geocoding.Location) v0.GeocodingOutput {
	data := make([]v0.PointOutput, len(locations))
	for i, location := range locations {
		data[i] = MapLocationToPointOutput(location)
	}

	return v0.GeocodingOutput{
		Count: len(locations),
		Data:  data,
	}
}
