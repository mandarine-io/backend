package e2e

import (
	"bytes"
	"github.com/goccy/go-json"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func ReadResponseBody(resp *http.Response, body interface{}) error {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Warn().Err(err).Msg("failed to close response body")
		}
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, body)
}

func NewJSONReader(body interface{}) (io.Reader, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(bodyBytes), nil
}

func MustMarshal(t provider.T, v interface{}) []byte {
	marshal, err := json.Marshal(v)
	t.Require().NoError(err)

	return marshal
}
