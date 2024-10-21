package api

import _ "embed"

var (
	//go:embed swagger.yaml
	SwaggerYaml []byte

	//go:embed swagger.json
	SwaggerJson []byte

	//go:embed ui/swagger-ui.css
	SwaggerUICSS []byte

	//go:embed ui/swagger-ui-bundle.js
	SwaggerUIBundleJS []byte

	//go:embed ui/swagger-ui-standalone-preset.js
	SwaggerUIStandalonePresetJS []byte
)
