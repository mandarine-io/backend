package geocoding

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"net/url"
)

type Provider interface {
	Geocode(ctx context.Context, address string, config GeocodeConfig) ([]*Location, error)
	ReverseGeocode(ctx context.Context, loc Location, config ReverseGeocodeConfig) ([]*Address, error)
}

type EndpointBuilder interface {
	GeocodeURL(address string, cfg GeocodeConfig) string
	ReverseGeocodeURL(loc Location, cfg ReverseGeocodeConfig) string
}

type ResponseParserFactory func() ResponseParser

type ResponseParser interface {
	Locations() ([]*Location, error)
	Addresses() ([]*Address, error)
	Error() error
}

type Option func(*baseProvider)

func WithLogger(logger zerolog.Logger) Option {
	return func(b *baseProvider) {
		b.logger = logger
	}
}

type baseProvider struct {
	endpointBuilder       EndpointBuilder
	responseParserFactory ResponseParserFactory

	logger zerolog.Logger
}

func NewBaseProvider(
	endpointBuilder EndpointBuilder,
	responseParserFactory ResponseParserFactory,
	opts ...Option,
) Provider {
	p := &baseProvider{
		endpointBuilder:       endpointBuilder,
		responseParserFactory: responseParserFactory,
		logger:                zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (g *baseProvider) Geocode(
	ctx context.Context,
	address string,
	cfg GeocodeConfig,
) ([]*Location, error) {
	g.logger.Debug().Msgf("geocode address: %s", address)

	responseParser := g.responseParserFactory()

	type geoResp struct {
		l []*Location
		e error
	}
	ch := make(chan geoResp, 1)

	go func(ch chan geoResp) {
		err := response(ctx, g.endpointBuilder.GeocodeURL(url.QueryEscape(address), cfg), responseParser, cfg.Lang)
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
		return nil, fmt.Errorf("failed to geocode: %w", ctx.Err())
	case res := <-ch:
		return res.l, fmt.Errorf("failed to geocode: %w", res.e)
	}
}

func (g *baseProvider) ReverseGeocode(
	ctx context.Context,
	loc Location,
	cfg ReverseGeocodeConfig,
) ([]*Address, error) {
	responseParser := g.responseParserFactory()

	type revResp struct {
		a []*Address
		e error
	}
	ch := make(chan revResp, 1)

	go func(ch chan revResp) {
		err := response(ctx, g.endpointBuilder.ReverseGeocodeURL(loc, cfg), responseParser, cfg.Lang)
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
		return nil, fmt.Errorf("failed to reverse geocode: %w", ctx.Err())
	case res := <-ch:
		return res.a, fmt.Errorf("failed to reverse geocode: %w", res.e)
	}
}

func response(ctx context.Context, url string, obj ResponseParser, lang language.Tag) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	req.Header.Add("User-Agent", "geo-golang/1.0")
	req.Header.Add("Accept-Language", lang.String())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, obj); err != nil {
		return err
	}

	if err = obj.Error(); err != nil {
		return err
	}

	return nil
}
