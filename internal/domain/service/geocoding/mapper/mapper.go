package mapper

import (
	"github.com/mandarine-io/Backend/internal/domain/dto"
	"github.com/mandarine-io/Backend/pkg/geocoding"
)

func MapAddressToAddressOutput(address *geocoding.Address) dto.AddressOutput {
	return dto.AddressOutput{
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

func MapAddressesToGeocodingOutput(addresses []*geocoding.Address) dto.ReverseGeocodingOutput {
	data := make([]dto.AddressOutput, len(addresses))
	for i, address := range addresses {
		data[i] = MapAddressToAddressOutput(address)
	}

	return dto.ReverseGeocodingOutput{
		Count: len(addresses),
		Data:  data,
	}
}

func MapLocationToPointOutput(location *geocoding.Location) dto.PointOutput {
	return dto.PointOutput{
		Latitude:  location.Lat,
		Longitude: location.Lng,
	}
}

func MapLocationsToGeocodingOutput(locations []*geocoding.Location) dto.GeocodingOutput {
	data := make([]dto.PointOutput, len(locations))
	for i, location := range locations {
		data[i] = MapLocationToPointOutput(location)
	}

	return dto.GeocodingOutput{
		Count: len(locations),
		Data:  data,
	}
}
