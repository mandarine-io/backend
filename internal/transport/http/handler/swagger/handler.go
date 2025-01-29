package swagger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apihandler "github.com/mandarine-io/backend/internal/transport/http/handler"
	"github.com/mandarine-io/backend/pkg/model/swagger"
	"github.com/rs/zerolog"
	"net/http"
)

type handler struct {
	swaggerYAML []byte
	swaggerJSON []byte
	uiStatic    []byte
	logger      zerolog.Logger
}

type Option func(*handler)

func WithLogger(logger zerolog.Logger) Option {
	return func(h *handler) {
		h.logger = logger
	}
}

func NewHandler(opts ...Option) apihandler.APIHandler {
	h := &handler{
		swaggerYAML: swagger.SwaggerYAML,
		swaggerJSON: swagger.SwaggerJSON,
		uiStatic:    renderUITemplate(),
		logger:      zerolog.Nop(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *handler) RegisterRoutes(router *gin.Engine) {
	h.logger.Debug().Msg("register swagger routes")

	swaggerRouter := router.Group("/swagger")
	{
		swaggerRouter.GET("/api-docs.json", h.getAPIDocJSON)
		swaggerRouter.GET("/api-docs.yaml", h.getAPIDocYAML)
		swaggerRouter.GET("/index.html", h.getUI)
	}
}

// getUI godoc
//
//	@Id				SwaggerUI
//	@Summary		Swagger UI
//	@Description	Request for getting swagger UI
//	@Tags			Swagger API
//	@Produce		text/html; charset=utf-8
//	@Success		200	{object}	string
//	@Router			/swagger/index.html [get]
func (h *handler) getUI(ctx *gin.Context) {
	h.logger.Debug().Msg("get swagger UI")
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", h.uiStatic)
}

// getAPIDocYAML godoc
//
//	@Id				SwaggerYAML
//	@Summary		Swagger YAML
//	@Description	Request for getting swagger specification in YAML
//	@Tags			Swagger API
//	@Produce		application/yaml; charset=utf-8
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.yaml [get]
func (h *handler) getAPIDocYAML(ctx *gin.Context) {
	h.logger.Debug().Msg("get swagger YAML")
	ctx.Data(http.StatusOK, "application/yaml; charset=utf-8", h.swaggerYAML)
}

// getAPIDocJSON godoc
//
//	@Id				SwaggerJSON
//	@Summary		Swagger JSON
//	@Description	Request for getting swagger specification in JSON
//	@Tags			Swagger API
//	@Produce		application/json; charset=utf-8
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.json [get]
func (h *handler) getAPIDocJSON(ctx *gin.Context) {
	h.logger.Debug().Msg("get swagger JSON")
	ctx.Data(http.StatusOK, "application/json; charset=utf-8", h.swaggerJSON)
}

func renderUITemplate() []byte {
	template := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.19.5/swagger-ui.css" >
    <style>
        .topbar {
            display: none;
        }
    </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.19.5/swagger-ui-bundle.js"> </script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.19.5/swagger-ui-standalone-preset.js"> </script>
<script>
	const spec = %s;
	window.onload = function() {
		// Build a system
		const ui = SwaggerUIBundle({
			dom_id: '#swagger-ui',
			deepLinking: true,
			spec: spec,
			presets: [
				SwaggerUIBundle.presets.apis,
				SwaggerUIStandalonePreset
			],
			plugins: [
				SwaggerUIBundle.plugins.DownloadURL
			],
			layout: "BaseLayout",
		})
		window.ui = ui
	}
</script>
</body>
</html>`

	return []byte(fmt.Sprintf(template, string(swagger.SwaggerJSON)))
}
