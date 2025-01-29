package yandex

import (
	"errors"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"github.com/shopspring/decimal"
	"strings"
)

const (
	componentTypeHouseNumber   = "house"
	componentTypeStreetName    = "street"
	componentTypeLocality      = "locality"
	componentTypeStateDistrict = "area"
	componentTypeState         = "province"
	componentTypeCountry       = "country"
)

type geocodeResponse struct {
	Response struct {
		GeoObjectCollection struct {
			MetaDataProperty struct {
				GeocoderResponseMetaData struct {
					Found string `json:"found"`
				} `json:"GeocoderResponseMetaData"`
			} `json:"metaDataProperty"`
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"`
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

type reverseGeocodeResponse struct {
	Response struct {
		GeoObjectCollection struct {
			MetaDataProperty struct {
				GeocoderResponseMetaData struct {
					Found string `json:"found"`
				} `json:"GeocoderResponseMetaData"`
			} `json:"metaDataProperty"`
			FeatureMember []struct {
				GeoObject struct {
					MetaDataProperty struct {
						GeocoderMetaData struct {
							Address struct {
								CountryCode string `json:"country_code"`
								PostalCode  string `json:"postal_code"`
								Formatted   string `json:"formatted"`
								Components  []struct {
									Kind string `json:"kind"`
									Name string `json:"name"`
								} `json:"Components"`
							} `json:"Address"`
						} `json:"GeocoderMetaData"`
					} `json:"metaDataProperty"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

type errorOutput struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type responseParser struct {
	geocodeResponse
	reverseGeocodeResponse
	errorOutput
}

func (r *geocodeResponse) Locations() ([]*geocoding.Location, error) {
	if r.Response.GeoObjectCollection.MetaDataProperty.GeocoderResponseMetaData.Found == "0" {
		return []*geocoding.Location{}, nil
	}
	if len(r.Response.GeoObjectCollection.FeatureMember) == 0 {
		return []*geocoding.Location{}, nil
	}

	result := make([]*geocoding.Location, len(r.Response.GeoObjectCollection.FeatureMember))
	for i, featureMember := range r.Response.GeoObjectCollection.FeatureMember {
		latLng := strings.Split(featureMember.GeoObject.Point.Pos, " ")
		if len(latLng) > 1 {
			lng, err := decimal.NewFromString(latLng[0])
			if err != nil {
				continue
			}

			lat, err := decimal.NewFromString(latLng[1])
			if err != nil {
				continue
			}

			result[i] = &geocoding.Location{
				Lat: lat,
				Lng: lng,
			}
		}
	}

	return result, nil
}

func (r *reverseGeocodeResponse) Addresses() ([]*geocoding.Address, error) {
	if r.Response.GeoObjectCollection.MetaDataProperty.GeocoderResponseMetaData.Found == "0" {
		return []*geocoding.Address{}, nil
	}
	if len(r.Response.GeoObjectCollection.FeatureMember) == 0 {
		return []*geocoding.Address{}, nil
	}

	addrs := make([]*geocoding.Address, len(r.Response.GeoObjectCollection.FeatureMember))
	for i, result := range r.Response.GeoObjectCollection.FeatureMember {
		res := result.GeoObject.MetaDataProperty.GeocoderMetaData
		addr := &geocoding.Address{}

		for _, comp := range res.Address.Components {
			switch comp.Kind {
			case componentTypeHouseNumber:
				addr.HouseNumber = comp.Name
				continue
			case componentTypeStreetName:
				addr.Street = comp.Name
				continue
			case componentTypeLocality:
				addr.City = comp.Name
				continue
			case componentTypeStateDistrict:
				addr.StateDistrict = comp.Name
				continue
			case componentTypeState:
				addr.State = comp.Name
				continue
			case componentTypeCountry:
				addr.Country = comp.Name
				continue
			}
		}

		addr.Postcode = res.Address.PostalCode
		addr.CountryCode = res.Address.CountryCode
		addr.FormattedAddress = res.Address.Formatted

		addrs[i] = addr
	}

	return addrs, nil
}

func (r *errorOutput) Error() error {
	return errors.New(r.Message)
}
