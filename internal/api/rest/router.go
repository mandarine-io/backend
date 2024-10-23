package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	docs "mandarine/docs/api"
	httphelper "mandarine/internal/api/helper/http"
	"mandarine/internal/api/registry"
	"mandarine/internal/api/rest/handler"
	"mandarine/pkg/rest/dto"
	"mandarine/pkg/rest/middleware"
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Setup method not allowed and route not found
	router.HandleMethodNotAllowed = true
	router.NoMethod(func(ctx *gin.Context) {
		_ = ctx.AbortWithError(http.StatusMethodNotAllowed, ErrMethodNotAllowed)
	})
	router.NoRoute(func(ctx *gin.Context) {
		_ = ctx.AbortWithError(http.StatusNotFound, ErrRouteNotFound)
	})

	// Setup middlewares
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.LocaleMiddleware(container.Bundle))
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.RateLimitMiddleware(container.RedisClient, container.Config.Security.RateLimit.RPS))
	router.Use(middleware.CorsMiddleware())

	if container.Config.Server.ExternalOrigin != "" && httphelper.IsPublicOrigin(container.Config.Server.ExternalOrigin) {
		router.Use(middleware.SecurityHeadersMiddleware())
	}

	// Register routes
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
		slog.Info(fmt.Sprintf("Register route: %6s %s", routeInfo.Method, routeInfo.Path))
	}

	return router
}
