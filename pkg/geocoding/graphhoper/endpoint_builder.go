package graphhopper

import (
	"fmt"
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"golang.org/x/text/language"
	"reflect"
)

type endpointBuilder struct {
	apiKey     string
	geocodeUrl string
	reverseUrl string
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

	return fmt.Sprintf("%s?q=%s&limit=%d&lang=%s&key=%s&reverse=false",
		b.geocodeUrl, address, limit, lang.String(), b.apiKey)
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

	return fmt.Sprintf("%s?point=%f,%f&limit=%d&lang=%s&key=%s&reverse=true",
		b.reverseUrl, l.Lat, l.Lng, limit, lang.String(), b.apiKey)
}
