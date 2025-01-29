package nominatim

import (
	"fmt"
	"github.com/mandarine-io/backend/third_party/geocoding"
)

type endpointBuilder struct {
	forGeocode, forReverseGeocode string
}

func (b *endpointBuilder) GeocodeURL(address string, cfg geocoding.GeocodeConfig) string {
	limit := 1
	if cfg.Limit > 0 {
		limit = cfg.Limit
	}
	return fmt.Sprintf("%s?format=json&q=%s&limit=%d", b.forGeocode, address, limit)
}

func (b *endpointBuilder) ReverseGeocodeURL(l geocoding.Location, cfg geocoding.ReverseGeocodeConfig) string {
	zoom := 18
	if cfg.Zoom > 0 {
		zoom = cfg.Zoom
	}
	return b.forReverseGeocode + fmt.Sprintf("?format=json&lat=%f&lon=%f&zoom=%d", l.Lat, l.Lng, zoom)
}
