package http

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/mandarine-io/Backend/pkg/geocoding"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/url"
	"time"
)

const DefaultTimeout = time.Second * 10

var ErrTimeout = errors.New("timeout")

type EndpointBuilder interface {
	GeocodeURL(address string, cfg geocoding.GeocodeConfig) string
	ReverseGeocodeURL(loc geocoding.Location, cfg geocoding.ReverseGeocodeConfig) string
}

type ResponseParserFactory func() ResponseParser

type ResponseParser interface {
	Locations() ([]*geocoding.Location, error)
	Addresses() ([]*geocoding.Address, error)
	Error() error
}

type Provider struct {
	EndpointBuilder       EndpointBuilder
	ResponseParserFactory ResponseParserFactory
}

func (g *Provider) GeocodeWithContext(
	ctx context.Context,
	address string,
	cfg geocoding.GeocodeConfig,
) ([]*geocoding.Location, error) {
	responseParser := g.ResponseParserFactory()

	type geoResp struct {
		l []*geocoding.Location
		e error
	}
	ch := make(chan geoResp, 1)

	go func(ch chan geoResp) {
		err := response(ctx, g.EndpointBuilder.GeocodeURL(url.QueryEscape(address), cfg), responseParser, cfg.Lang)
		if err != nil {
			ch <- geoResp{
				l: nil,
				e: err,
			}
		}

		loc, err := responseParser.Locations()
		ch <- geoResp{
			l: loc,
			e: err,
		}
	}(ch)

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case res := <-ch:
		return res.l, res.e
	}
}

func (g *Provider) Geocode(address string, cfg geocoding.GeocodeConfig) ([]*geocoding.Location, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), DefaultTimeout)
	defer cancel()

	return g.GeocodeWithContext(ctx, address, cfg)
}

func (g *Provider) ReverseGeocodeWithContext(
	ctx context.Context,
	loc geocoding.Location,
	cfg geocoding.ReverseGeocodeConfig,
) ([]*geocoding.Address, error) {
	responseParser := g.ResponseParserFactory()

	type revResp struct {
		a []*geocoding.Address
		e error
	}
	ch := make(chan revResp, 1)

	go func(ch chan revResp) {
		err := response(ctx, g.EndpointBuilder.ReverseGeocodeURL(loc, cfg), responseParser, cfg.Lang)
		if err != nil {
			ch <- revResp{
				a: nil,
				e: err,
			}
		}

		addr, err := responseParser.Addresses()
		ch <- revResp{
			a: addr,
			e: err,
		}
	}(ch)

	select {
	case <-ctx.Done():
		return nil, ErrTimeout
	case res := <-ch:
		return res.a, res.e
	}
}

func (g *Provider) ReverseGeocode(
	loc geocoding.Location,
	cfg geocoding.ReverseGeocodeConfig,
) ([]*geocoding.Address, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), DefaultTimeout)
	defer cancel()

	return g.ReverseGeocodeWithContext(ctx, loc, cfg)
}

func response(ctx context.Context, url string, obj ResponseParser, lang language.Tag) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	req.Header.Add("User-Agent", "geo-golang/1.0")
	req.Header.Add("Accept-Language", lang.String())

	log.Debug().Msgf("send request to %s", req.URL.String())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn().Stack().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debug().Msgf("response body: %s", string(body))

	if err := json.Unmarshal(body, obj); err != nil {
		log.Error().Stack().Err(err).Msg("failed to unmarshal response")
		return err
	}

	if err := obj.Error(); err != nil {
		log.Error().Stack().Err(err).Msg("failed to get error from response")
		return geocoding.NewProviderError(err.Error(), resp.StatusCode)
	}

	return nil
}
