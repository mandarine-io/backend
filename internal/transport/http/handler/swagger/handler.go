package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/docs/api"
	apihandler "github.com/mandarine-io/Backend/internal/transport/http/handler"
	"github.com/rs/zerolog/log"
	"net/http"
)

type handler struct {
	swaggerYaml []byte
	swaggerJson []byte
	uiStatic    []byte
}

func NewHandler() apihandler.ApiHandler {
	return &handler{
		swaggerYaml: api.SwaggerYaml,
		swaggerJson: api.SwaggerJson,
		uiStatic:    renderUITemplate(),
	}
}

func (h *handler) RegisterRoutes(router *gin.Engine, _ apihandler.RouteMiddlewares) {
	log.Debug().Msg("register swagger routes")

	router.GET("/swagger/api-docs.json", h.GetApiDocJson)
	router.GET("/swagger/api-docs.yaml", h.GetApiDocYaml)
	router.GET("/swagger/index.html", h.GetUI)
}

// GetUI godoc
//
//	@Id				SwaggerUI
//	@Summary		Swagger UI
//	@Description	Request for getting swagger UI
//	@Tags			Swagger API
//	@Produce		text/html
//	@Success		200	{object}	string
//	@Router			/swagger/index.html [get]
func (h *handler) GetUI(ctx *gin.Context) {
	log.Debug().Msg("get swagger ui")
	ctx.Data(http.StatusOK, "text/html", h.uiStatic)
}

// GetApiDocYaml godoc
//
//	@Id				Swagger API specification in YAML
//	@Summary		Swagger YAML
//	@Description	Request for getting swagger specification in YAML
//	@Tags			Swagger API
//	@Produce		text/plain
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.yaml [get]
func (h *handler) GetApiDocYaml(ctx *gin.Context) {
	log.Debug().Msg("get swagger yaml")
	ctx.Data(http.StatusOK, "text/plain", h.swaggerYaml)
}

// GetApiDocJson godoc
//
//	@Id				Swagger API specification in JSON
//	@Summary		Swagger JSON
//	@Description	Request for getting swagger specification in JSON
//	@Tags			Swagger API
//	@Produce		text/plain
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.json [get]
func (h *handler) GetApiDocJson(ctx *gin.Context) {
	log.Debug().Msg("get swagger json")
	ctx.Data(http.StatusOK, "text/plain", h.swaggerJson)
}

func renderUITemplate() []byte {
	return []byte(`<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta charset="UTF-8">
	<title>Mandarine API</title>
<style>
` + string(api.SwaggerUICSS) + `
</style>
</head>
<body>

<div id="swagger-ui"></div>

<script>` + string(api.SwaggerUIBundleJS) + `</script>
<script>` + string(api.SwaggerUIStandalonePresetJS) + `</script>

<script>
	const spec = ` + string(api.SwaggerJson) + `;
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
				SwaggerUIBundle.plugins.DownloadUrl
			],
			layout: "BaseLayout",
		})
		window.ui = ui
	}
</script>
</body>
</html>`)
}
