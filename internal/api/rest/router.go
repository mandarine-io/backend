package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	docs "mandarine/docs/api"
	httphelper "mandarine/internal/api/helper/http"
	"mandarine/internal/api/registry"
	"mandarine/internal/api/rest/handler/common"
	"mandarine/pkg/logging"
	"mandarine/pkg/rest/middleware"
	"sort"
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
//	@tag.name					Metrics API
//	@tag.description			API for getting metrics
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

	// Setup middlewares
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.LocaleMiddleware(container.Bundle))
	router.Use(middleware.ProduceMiddleware())
	router.Use(middleware.ErrorMiddleware())
	router.Use(middleware.RateLimitMiddleware(container.RedisClient, container.Config.Security.RateLimit.RPS))
	router.Use(middleware.CorsMiddleware())

	if container.Config.Server.ExternalOrigin != "" && httphelper.IsPublicOrigin(container.Config.Server.ExternalOrigin) {
		router.Use(middleware.SecurityHeadersMiddleware())
	}

	requireAuth := middleware.JWTMiddleware(container.Config.Security.JWT, container.Repositories.BannedToken)

	// Common
	router.HandleMethodNotAllowed = true
	router.NoMethod(common.NoMethod)
	router.NoRoute(common.NoRoute)

	// Health Check
	err := common.SetupHealthcheck(router, container.DB, container.RedisClient, container.S3, container.SmtpSender)
	if err != nil {
		slog.Warn("Health check route error", logging.ErrorAttr(err))
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// V0 routes
	v0Handlers := container.Handlers.V0
	v0Router := router.Group("v0")
	{
		authRouter := v0Router.Group("auth")
		{
			authRouter.POST("login", v0Handlers.Login.Login)
			authRouter.GET("refresh", v0Handlers.Login.RefreshTokens)
			authRouter.GET("social/:provider", v0Handlers.SocialLogin.SocialLogin)
			authRouter.POST("social/:provider/callback", v0Handlers.SocialLogin.SocialLoginCallback)
			authRouter.GET("logout", requireAuth, v0Handlers.Logout.Logout)
			authRouter.POST("register", v0Handlers.Register.Register)
			authRouter.POST("register/confirm", v0Handlers.Register.RegisterConfirm)
			authRouter.POST("recovery-password", v0Handlers.ResetPassword.RecoveryPassword)
			authRouter.POST("recovery-password/verify", v0Handlers.ResetPassword.VerifyRecoveryCode)
			authRouter.POST("reset-password", v0Handlers.ResetPassword.ResetPassword)
		}

		accountRouter := v0Router.Group("account")
		{
			accountRouter.GET("", requireAuth, v0Handlers.Account.GetAccount)
			accountRouter.PATCH("username", requireAuth, v0Handlers.Account.UpdateUsername)
			accountRouter.PATCH("email", requireAuth, v0Handlers.Account.UpdateEmail)
			accountRouter.POST("email/verify", requireAuth, v0Handlers.Account.VerifyEmail)
			accountRouter.POST("password", requireAuth, v0Handlers.Account.SetPassword)
			accountRouter.PATCH("password", requireAuth, v0Handlers.Account.UpdatePassword)
			accountRouter.DELETE("", requireAuth, v0Handlers.Account.DeleteAccount)
			accountRouter.GET("restore", requireAuth, v0Handlers.Account.RestoreAccount)
		}
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
