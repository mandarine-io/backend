package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/mandarine-io/backend/internal/infrastructure/websocket"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

const (
	poolSize = 10
)

var (
	pool   *websocket.Pool
	server *httptest.Server
)

type WebsocketPoolSuite struct {
	suite.Suite
}

func TestWebsocketPoolSuite(t *testing.T) {
	var err error
	pool, err = websocket.NewPool(poolSize)
	require.NoError(t, err)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET(
		"/ws/:id", func(c *gin.Context) {
			_ = pool.Register(c.Param("id"), c.Request, c.Writer)
		},
	)

	server = httptest.NewServer(router)
	defer server.Close()

	suite.RunSuite(t, new(WebsocketPoolSuite))
}

func (s *WebsocketPoolSuite) Test(t provider.T) {
	s.RunSuite(t, new(RegisterSendReceiveSuite))
}
