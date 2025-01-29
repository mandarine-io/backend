package here

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

	return fmt.Sprintf("%s?apiKey=%s&lang=%s&limit=%d&q=%s", b.geocodeUrl, b.apiKey, lang.String(), limit, address)
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
		"%s?apiKey=%s&lang=%s&limit=%d&at=%f,%f&types=city,street",
		b.reverseGeocodeUrl, b.apiKey, lang.String(), limit, l.Lat, l.Lng,
	)
}
