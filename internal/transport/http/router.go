package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/mandarine-io/backend/internal/di"
	"github.com/mandarine-io/backend/internal/transport/http/middleware"
	"github.com/mandarine-io/backend/internal/transport/http/util"
	validator2 "github.com/mandarine-io/backend/internal/transport/http/validator"
	"github.com/mandarine-io/backend/pkg/model/swagger"
	"github.com/mandarine-io/backend/pkg/model/v0"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"net/http"
	"sort"
)

var (
	ErrMethodNotAllowed = v0.NewI18nError("method not allowed", "errors.method_not_allowed")
	ErrRouteNotFound    = v0.NewI18nError("route not found", "errors.route_not_found")

	ignorePathRegexps = []string{
		"/metrics.*",
		"/health.*",
		"/swagger.*",
	}
)

type RequireRoleMiddlewareFactory func(...string) gin.HandlerFunc

// SetupRouter godoc
//
//	@title						Mandarine API
//	@version					0.0.0
//	@description				API for web and mobile application Mandarine
//	@host						localhost:8080
//	@query.collection.format	multi
//	@accept						json
//	@produce					json
//	@tag.name					Account API
//	@tag.description			API for account management
//	@tag.name					Authentication and Authorization API
//	@tag.description			API for authentication and authorization
//	@tag.name					Geocoding API
//	@tag.description			API for geocoding
//	@tag.name					Master Profile API
//	@tag.description			API for master profile management
//	@tag.name					Master Service API
//	@tag.description			API for master service management
//	@tag.name					Metrics API
//	@tag.description			API for getting metrics and healthcheck
//	@tag.name					Resource API
//	@tag.description			API for download and upload files
//	@tag.name					Swagger API
//	@tag.description			API for getting swagger documentation
//	@tag.name					Websocket API
//	@tag.description			API for establishing websocket connection
//	@contact.name				Mandarine Support
//	@contact.email				mandarine.app@yandex.ru
//	@license.name				Apache 2.0
//	@license.url				https://www.apache.org/licenses/LICENSE-2.0.html
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func SetupRouter(container *di.Container) *gin.Engine {
	// Setup Swagger spec
	swagger.SwaggerInfo.Version = container.Config.Server.Version
	swagger.SwaggerInfo.Host = container.Config.Server.ExternalURL

	// Create router
	log.Debug().Msg("create router")
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Setup validators
	log.Debug().Msg("setup validators")
	decimal.MarshalJSONWithoutQuotes = true
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("pastdate", validator2.PastDateValidator)
		_ = v.RegisterValidation("duration", validator2.DurationValidator)
		_ = v.RegisterValidation("zxcvbn", validator2.ZxcvbnPasswordValidator)
		_ = v.RegisterValidation("username", validator2.UsernameValidator)
		_ = v.RegisterValidation("point", validator2.PointValidator)
	}

	// Setup method not allowed and route not found
	log.Debug().Msg("setup method not allowed and route not found")
	router.HandleMethodNotAllowed = true
	router.NoMethod(handleNoMethod)
	router.NoRoute(handleNoRoute)

	// Setup middlewares
	log.Debug().Msg("setup middlewares")
	router.Use(middleware.ObservabilityMiddleware(container.Metrics, ignorePathRegexps...))
	router.Use(middleware.LoggerMiddleware(ignorePathRegexps...))
	router.Use(middleware.LocaleMiddleware(container.Infrastructure.LocaleBundle))
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.ErrorMiddleware())
	middleware.InitRegistry(container.InfrastructureSVCs.JWT)

	// Register routes
	log.Debug().Msg("register routes")
	for _, apiHandler := range container.Handlers {
		apiHandler.RegisterRoutes(router)
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

func handleNoMethod(ctx *gin.Context) {
	log.Debug().Msg("handle method not allowed")
	_ = util.ErrorWithStatus(ctx, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
}

func handleNoRoute(ctx *gin.Context) {
	log.Debug().Msg("handle route not found")
	_ = util.ErrorWithStatus(ctx, http.StatusNotFound, ErrRouteNotFound)
}
