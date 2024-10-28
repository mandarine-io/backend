package swagger

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/Backend/docs/api"
	"github.com/mandarine-io/Backend/internal/api/rest/handler"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Handler struct {
	swaggerYaml []byte
	swaggerJson []byte
	uiStatic    []byte
}

func NewHandler() *Handler {
	return &Handler{
		swaggerYaml: api.SwaggerYaml,
		swaggerJson: api.SwaggerJson,
		uiStatic:    renderUITemplate(),
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine, _ handler.RouteMiddlewares) {
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
func (h *Handler) GetUI(ctx *gin.Context) {
	log.Debug().Msg("get swagger ui")
	ctx.Data(http.StatusOK, "text/html", h.uiStatic)
}

// GetApiDocYaml godoc
//
//	@Id				Swagger API specification in YAML
//	@Summary		Swagger YAML
//	@Description	Request for getting swagger specification in YAML
//	@Tags			Swagger API
//	@Produce		application/yaml
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.yaml [get]
func (h *Handler) GetApiDocYaml(ctx *gin.Context) {
	log.Debug().Msg("get swagger yaml")
	ctx.Data(http.StatusOK, "application/yaml", h.swaggerYaml)
}

// GetApiDocJson godoc
//
//	@Id				Swagger API specification in JSON
//	@Summary		Swagger JSON
//	@Description	Request for getting swagger specification in JSON
//	@Tags			Swagger API
//	@Produce		application/yaml
//	@Success		200	{object}	string
//	@Router			/swagger/api-docs.json [get]
func (h *Handler) GetApiDocJson(ctx *gin.Context) {
	log.Debug().Msg("get swagger json")
	ctx.Data(http.StatusOK, "application/json", h.swaggerJson)
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
