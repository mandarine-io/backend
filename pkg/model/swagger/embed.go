package swagger

import _ "embed"

var (
	//go:embed swagger.yaml
	SwaggerYAML []byte

	//go:embed swagger.json
	SwaggerJSON []byte
)
