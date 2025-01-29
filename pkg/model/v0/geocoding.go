package v0

import "github.com/shopspring/decimal"

type GeocodingInput struct {
	Address string `form:"address" binding:"required"`
	Limit   int    `form:"limit,default=1" binding:"min=1"`
}

type GeocodingOutput struct {
	Count int           `json:"count"`
	Data  []PointOutput `json:"data"`
}

type ReverseGeocodingInput struct {
	Lng   decimal.Decimal `form:"lng" format:"decimal" binding:"required"`
	Lat   decimal.Decimal `form:"lat" format:"decimal" binding:"required"`
	Limit int             `form:"limit,default=1" binding:"min=1"`
}

type AddressOutput struct {
	FormattedAddress string `json:"formattedAddress"`
	Street           string `json:"street"`
	HouseNumber      string `json:"houseNumber"`
	Suburb           string `json:"suburb"`
	Postcode         string `json:"postcode"`
	State            string `json:"state"`
	StateCode        string `json:"stateCode"`
	StateDistrict    string `json:"stateDistrict"`
	County           string `json:"county"`
	Country          string `json:"country"`
	CountryCode      string `json:"countryCode"`
	City             string `json:"city"`
}

type ReverseGeocodingOutput struct {
	Count int             `json:"count"`
	Data  []AddressOutput `json:"data"`
}
