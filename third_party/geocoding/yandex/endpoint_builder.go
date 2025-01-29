package yandex

import (
	"fmt"
	"github.com/mandarine-io/backend/third_party/geocoding"
	"golang.org/x/text/language"
	"reflect"
)

type endpointBuilder struct {
	apiKey            string
	geocodeUrl        string
	reverseGeocodeUrl string
}

func (b *endpointBuilder) GeocodeURL(address string, cfg geocoding.GeocodeConfig) string {
	limit := 1
	if cfg.Limit > 0 {
		limit = cfg.Limit
	}

	lang := language.English
	if !reflect.ValueOf(cfg.Lang).IsZero() {
		lang = cfg.Lang
	}

	return fmt.Sprintf(
		"%s?results=%d&lang=%s&format=json&apikey=%s&geocode=%s",
		b.geocodeUrl, limit, lang.String(), b.apiKey, address,
	)
}

func (b *endpointBuilder) ReverseGeocodeURL(l geocoding.Location, cfg geocoding.ReverseGeocodeConfig) string {
	limit := 1
	if cfg.Limit > 0 {
		limit = cfg.Limit
	}

	lang := language.English
	if !reflect.ValueOf(cfg.Lang).IsZero() {
		lang = cfg.Lang
	}

	return fmt.Sprintf(
		"%s?results=%d&lang=%s&format=json&kind=house&apikey=%s&sco=latlong&geocode=%f,%f",
		b.geocodeUrl, limit, lang.String(), b.apiKey, l.Lat, l.Lng,
	)
}
