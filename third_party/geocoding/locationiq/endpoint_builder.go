package locationiq

import (
	"fmt"
	"github.com/mandarine-io/backend/third_party/geocoding"
)

type endpointBuilder struct {
	geocodeUrl        string
	reverseGeocodeUrl string
	apiKey            string
}

func (b *endpointBuilder) GeocodeURL(address string, cfg geocoding.GeocodeConfig) string {
	limit := 1
	if cfg.Limit > 0 {
		limit = cfg.Limit
	}
	return fmt.Sprintf("%s?key=%s&format=json&limit=%d&q=%s", b.geocodeUrl, b.apiKey, limit, address)
}

func (b *endpointBuilder) ReverseGeocodeURL(l geocoding.Location, cfg geocoding.ReverseGeocodeConfig) string {
	zoom := 18
	if cfg.Zoom > 0 {
		zoom = cfg.Zoom
	}
	return fmt.Sprintf("%s?key=%s&format=json&lat=%f&lon=%f&zoom=%d", b.reverseGeocodeUrl, b.apiKey, l.Lat, l.Lng, zoom)
}
