package smtp

import (
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type SendPlainMessageSuite struct {
	suite.Suite
}

func (s *SendPlainMessageSuite) Test_Success(t provider.T) {
	t.Title("Send plain message - success")
	t.Severity(allure.NORMAL)
	t.Feature("SMTP sender")
	t.Tags("Positive")
	t.Parallel()

	err := sender.SendPlainMessage("subject", "content", "sender_plain@mail.ru", "receiver@mail.ru")
	t.Require().NoError(err)

	client := resty.New()

	resp, err := client.R().Get(mailhogApiURL + "/api/v2/messages")
	t.Require().NoError(err)
	t.Require().Equal(200, resp.StatusCode())

	messages := GetMessagesResponse{}
	err = json.Unmarshal(resp.Body(), &messages)
	t.Require().NoError(err)

	for _, message := range messages.Items {
		if message.Raw.From == "sender_plain@mail.ru" && message.Raw.To[0] == "receiver@mail.ru" {
			t.Require().Equal("subject", message.Content.Headers.Subject[0])
			t.Require().Equal("content", message.Content.Body)
			return
		}
	}

	t.Fail()
}

func (s *SendPlainMessageSuite) Test_IncorrectEmail(t provider.T) {
	t.Title("Send plain message - incorrect email")
	t.Severity(allure.CRITICAL)
	t.Feature("SMTP sender")
	t.Tags("Negative")
	t.Parallel()

	err := sender.SendPlainMessage("subject", "content", "sender_plain_incorrect_email@mail.ru", "not_email")
	t.Require().Error(err)
	t.Require().Contains(err.Error(), "invalid address")
}
