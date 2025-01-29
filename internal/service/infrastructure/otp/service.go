package otp

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/mandarine-io/backend/config"
	"github.com/mandarine-io/backend/internal/infrastructure/cache"
	"github.com/mandarine-io/backend/internal/service/infrastructure"
	cachehelper "github.com/mandarine-io/backend/internal/util/cache"
	"github.com/rs/zerolog"
	"math/big"
	"time"
)

var (
	ErrNegativeOTPLength = errors.New("OTP length is negative")
)

type svc struct {
	manager cache.Manager
	cfg     config.OTPConfig
	logger  zerolog.Logger
}

type Option func(*svc)

func WithLogger(logger zerolog.Logger) Option {
	return func(p *svc) {
		p.logger = logger
	}
}

func NewService(manager cache.Manager, cfg config.OTPConfig, opts ...Option) infrastructure.OTPService {
	p := &svc{
		manager: manager,
		cfg:     cfg,
		logger:  zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (s *svc) GenerateCode(_ context.Context) (string, error) {
	s.logger.Debug().Msg("create OTP code")
	return generateRandomNumber(s.cfg.Length)
}

func (s *svc) SaveWithCode(ctx context.Context, prefix string, code string, data any) error {
	s.logger.Debug().Msg("save OTP code")

	return s.manager.SetWithExpiration(
		ctx,
		cachehelper.CreateCacheKey(prefix, code),
		data,
		time.Duration(s.cfg.TTL)*time.Second,
	)
}

func (s *svc) GenerateAndSaveWithCode(ctx context.Context, prefix string, data any) (string, error) {
	code, err := s.GenerateCode(ctx)
	if err != nil {
		return "", err
	}

	err = s.SaveWithCode(ctx, prefix, code, data)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (s *svc) GetDataByCode(ctx context.Context, prefix string, code string, data any) error {
	s.logger.Debug().Msg("get OTP code")

	err := s.manager.Get(ctx, cachehelper.CreateCacheKey(prefix, code), data)
	if errors.Is(err, cache.ErrCacheEntryNotFound) {
		return infrastructure.ErrInvalidOrExpiredOTP
	}

	return err
}

func (s *svc) DeleteDataByCode(ctx context.Context, prefix string, code string) error {
	s.logger.Debug().Msg("delete OTP code")
	return s.manager.Delete(ctx, cachehelper.CreateCacheKey(prefix, code))
}

func generateRandomNumber(length int) (string, error) {
	if length < 0 {
		return "", ErrNegativeOTPLength
	}

	const digits = "0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		if num == nil {
			num = big.NewInt(0)
		}

		result[i] = digits[num.Int64()]
	}

	return string(result), nil
}
