package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	websocket2 "github.com/mandarine-io/backend/internal/infrastructure/websocket"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"net/url"
	"strings"
)

type RegisterSendReceiveSuite struct {
	suite.Suite
}

func (suite *RegisterSendReceiveSuite) Test_Success(t provider.T) {
	t.Title("Register send and receive - success")
	t.Severity(allure.NORMAL)
	t.Feature("Websocket pool")
	t.Tags("Positive")
	t.Parallel()

	pool.RegisterHandler(
		func(pool *websocket2.Pool, msg websocket2.ClientMessage) {
			pool.Send(msg.ClientID, msg.Payload)
		},
	)

	u := url.URL{
		Scheme: "ws",
		Host:   strings.Replace(server.URL, "http://", "", 1),
		Path:   "/ws/1",
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	t.Require().NoError(err)
	defer func(c *websocket.Conn) {
		_ = c.Close()
	}(c)

	err = c.WriteMessage(websocket.TextMessage, []byte("Hello"))
	t.Require().NoError(err)

	_, message, err := c.ReadMessage()
	t.Require().NoError(err)
	t.Require().Equal("Hello", string(message))
}

func (suite *RegisterSendReceiveSuite) Test_PoolIsFull(t provider.T) {
	t.Title("Register send and receive - pool is full")
	t.Severity(allure.CRITICAL)
	t.Feature("Websocket pool")
	t.Tags("Negative")

	conns := make([]*websocket.Conn, poolSize)
	for i := 0; i < poolSize; i++ {
		u := url.URL{
			Scheme: "ws",
			Host:   strings.Replace(server.URL, "http://", "", 1),
			Path:   fmt.Sprintf("/ws/%d", i),
		}
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		t.Require().NoError(err)
		conns[i] = c
	}
	defer func() {
		for _, c := range conns {
			_ = c.Close()
		}
	}()

	u := url.URL{
		Scheme: "ws",
		Host:   strings.Replace(server.URL, "http://", "", 1),
		Path:   fmt.Sprintf("/ws/%d", poolSize),
	}
	_, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	t.Require().Error(err)
}
