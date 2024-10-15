package common_e2e_test

import (
	"github.com/stretchr/testify/assert"
	"mandarine/internal/api/rest/handler/common/dto"
	"mandarine/tests/e2e"
	"net/http"
	"testing"
)

func Test_HealthCheck(t *testing.T) {
	url := server.URL + "/health"

	resp, err := server.Client().Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	//// Check response
	var body []dto.HealthResponse
	err = e2e.ReadResponseBody(resp, &body)
	assert.NoError(t, err)
	for _, v := range body {
		assert.True(t, v.Pass)
	}
}
