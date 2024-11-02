package rest

import (
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	docs "github.com/mandarine-io/Backend/docs/api"
	httphelper "github.com/mandarine-io/Backend/internal/api/helper/http"
	"github.com/mandarine-io/Backend/internal/api/registry"
	"github.com/mandarine-io/Backend/internal/api/rest/handler"
	"github.com/mandarine-io/Backend/pkg/rest/dto"
	"github.com/mandarine-io/Backend/pkg/rest/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"sort"
)

var (
	ErrMethodNotAllowed = dto.NewI18nError("method not allowed", "errors.method_not_allowed")
	ErrRouteNotFound    = dto.NewI18nError("route not found", "errors.route_not_found")
)

type RequireRoleMiddlewareFactory func(...string) gin.HandlerFunc

// SetupRouter godoc
//
//	@title						Mandarine API
//	@version					0.0.0
//	@description				API for web and mobile application Mandarine
//	@host						localhost:8080
//	@accept						json
//	@produce					json
//	@tag.name					Account API
//	@tag.description			API for account management
//	@tag.name					Authentication and Authorization API
//	@tag.description			API for authentication and authorization
//	@tag.name					Resource API
//	@tag.description			API for resource management
//	@tag.name					Websocket API
//	@tag.description			API for websocket connection
//	@tag.name					Metrics API
//	@tag.description			API for getting metrics
//	@tag.name					Swagger API
//	@tag.description			API for getting swagger documentation
//	@contact.name				Mandarine Support
//	@contact.email				mandarine.app@yandex.ru
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func SetupRouter(container *registry.Container) *gin.Engine {
	docs.SwaggerInfo.Version = container.Config.Server.Version

	// Create router
	log.Debug().Msg("create router")
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Setup method not allowed and route not found
	log.Debug().Msg("setup method not allowed and route not found")
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(ctx *gin.Context) {
		log.Debug().Msg("handle method not allowed")
		_ = ctx.AbortWithError(http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	})
	router.NoRoute(func(ctx *gin.Context) {
		log.Debug().Msg("handle route not found")
		_ = ctx.AbortWithError(http.StatusNotFound, ErrRouteNotFound)
	})

	// Setup middlewares
	log.Debug().Msg("setup middlewares")
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.LocaleMiddleware(container.Bundle))
	router.Use(middleware.RateLimitMiddleware(container.RedisClient, container.Config.Security.RateLimit.RPS))
	router.Use(middleware.CorsMiddleware())
	router.Use(limits.RequestSizeLimiter(int64(container.Config.Server.MaxRequestSize)))
	router.Use(middleware.ErrorMiddleware())

	if container.Config.Server.ExternalOrigin != "" && httphelper.IsPublicOrigin(container.Config.Server.ExternalOrigin) {
		router.Use(middleware.SecurityHeadersMiddleware())
	}

	// Register routes
	log.Debug().Msg("register routes")
	middlewares := handler.RouteMiddlewares{
		Auth:        middleware.JWTMiddleware(container.Config.Security.JWT, container.Repositories.BannedToken),
		RoleFactory: middleware.RoleMiddleware,
		BannedUser:  middleware.BannedUserMiddleware(),
		DeletedUser: middleware.DeletedUserMiddleware(),
	}
	for _, apiHandler := range container.Handlers {
		apiHandler.RegisterRoutes(router, middlewares)
	}

	// Log routes
	routes := router.Routes()
	sort.Slice(
		routes, func(i, j int) bool {
			return routes[i].Path < routes[j].Path
		},
	)
	for _, routeInfo := range routes {
		log.Info().Msgf("registered route: %6s %s", routeInfo.Method, routeInfo.Path)
	}

	return router
}
