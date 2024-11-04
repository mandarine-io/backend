package middleware

import (
	"fmt"
	"github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/pkg/transport/http/dto"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

var (
	ErrTooManyRequests = dto.NewI18nError("too many requests", "errors.too_many_requests")
)

func MemoryRateLimitMiddleware(rps int) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Second,
			Limit: uint(rps),
		},
	)

	return rateLimitMiddleware(store)
}

func RedisRateLimitMiddleware(redisClient *redis.Client, rps int) gin.HandlerFunc {
	store := ratelimit.RedisStore(
		&ratelimit.RedisOptions{
			RedisClient: redisClient,
			Rate:        time.Second,
			Limit:       uint(rps),
		},
	)

	return rateLimitMiddleware(store)
}

func rateLimitMiddleware(store ratelimit.Store) gin.HandlerFunc {
	log.Debug().Msg("setup rate limit middleware")

	keyFunc := func(c *gin.Context) string {
		return c.ClientIP()
	}

	errorHandler := func(c *gin.Context, info ratelimit.Info) {
		log.Debug().Msg("set rate limit headers")

		c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", info.Limit))
		c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", info.ResetTime.Unix()))
		_ = c.AbortWithError(http.StatusTooManyRequests, ErrTooManyRequests)
	}

	beforeResponse := func(c *gin.Context, info ratelimit.Info) {
		log.Debug().Msg("set rate limit headers")

		c.Header("X-Rate-Limit-Limit", fmt.Sprintf("%d", info.Limit))
		c.Header("X-Rate-Limit-Remaining", fmt.Sprintf("%v", info.RemainingHits))
		c.Header("X-Rate-Limit-Reset", fmt.Sprintf("%d", info.ResetTime.Unix()))
	}

	return ratelimit.RateLimiter(
		store, &ratelimit.Options{
			KeyFunc:        keyFunc,
			BeforeResponse: beforeResponse,
			ErrorHandler:   errorHandler,
		},
	)
}
