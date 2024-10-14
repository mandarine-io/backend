package auth_e2e_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mandarine/internal/api/helper/security"
	"mandarine/internal/api/persistence/model"
	"mandarine/tests/e2e"
	"net/http"
	"testing"
)

func Test_LogoutHandler_Logout(t *testing.T) {
	e2e.MustLoadFixtures(fixtures)

	url := server.URL + "/v0/auth/logout"

	t.Run("Unauthorized", func(t *testing.T) {
		// Send request
		req, _ := http.NewRequest("GET", url, nil)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Success", func(t *testing.T) {
		// New access token
		accessToken, _, _ := security.GenerateTokens(testEnvironment.Container.Config.Security.JWT, &model.UserEntity{
			ID:       uuid.MustParse("a83d9587-b01f-4146-8b1f-80f137f53534"),
			Username: "user_for_logout",
			Email:    "user_for_logout@example.com",
			Role:     model.RoleEntity{Name: model.RoleUser},
		})

		// Send request
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err := server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check repeated request
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)

		resp, err = server.Client().Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
