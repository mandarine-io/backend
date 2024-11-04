package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func ReadResponseBody(resp *http.Response, body interface{}) error {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Warn().Stack().Err(err).Msg("failed to close response body")
		}
	}()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, body)
}

func NewJSONReader(body interface{}) (io.Reader, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}
