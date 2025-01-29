package auth

import (
	redis2 "github.com/mandarine-io/backend/internal/infrastructure/cache/redis"
	postgres2 "github.com/mandarine-io/backend/internal/infrastructure/database/gorm/postgres"
	"github.com/mandarine-io/backend/tests/e2e"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"testing"
)

var (
	secret     string
	mailhogURL string
	serverURL  string
	db         *gorm.DB
	rdb        redis.UniversalClient
)

func init() {
	serverURL = e2e.Cfg.GetServerURL()
	secret = e2e.Cfg.GetServerJWTSecret()

	var err error
	db, err = postgres2.NewDb(
		e2e.Cfg.GetPostgresConfig(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	rdb, err = redis2.NewClient(
		e2e.Cfg.GetRedisConfig(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}

	mailhogURL = e2e.Cfg.GetMailhogAPIURL()
}

type AuthHandlerSuite struct {
	suite.Suite
}

func TestAuthHandlerSuite(t *testing.T) {
	defer func(db *gorm.DB) {
		_ = postgres2.CloseDb(db)
	}(db)

	defer func(rdb redis.UniversalClient) {
		_ = rdb.Close()
	}(rdb)

	suite.RunSuite(t, new(AuthHandlerSuite))
}

func (s *AuthHandlerSuite) Test(t provider.T) {
	//s.RunSuite(t, new(LoginSuite))
	//s.RunSuite(t, new(LogoutSuite))
	//s.RunSuite(t, new(RefreshTokensSuite))
	//s.RunSuite(t, new(RegisterSuite))
	//s.RunSuite(t, new(ResetPasswordSuite))
}
