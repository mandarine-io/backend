package auth

import (
	"github.com/gin-gonic/gin"
	"mandarine/internal/api/service/auth"
	"mandarine/pkg/rest/middleware"
	"net/http"
)

type LogoutHandler struct {
	logoutService *auth.LogoutService
}

func NewLogoutHandler(logoutService *auth.LogoutService) LogoutHandler {
	return LogoutHandler{
		logoutService: logoutService,
	}
}

// Logout godoc
//
//	@Id				Logout
//	@Summary		Logout
//	@Description	Request for logout. User must be logged in.
//	@Security		BearerAuth
//	@Tags			Authentication and Authorization API
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/v0/auth/logout [get]
func (h *LogoutHandler) Logout(c *gin.Context) {
	principal, err := middleware.GetAuthUser(c)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	err = h.logoutService.Logout(c, principal.JTI)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
