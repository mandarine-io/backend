package smtp

import (
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type SendHTMLMessagesSuite struct {
	suite.Suite
}

func (s *SendHTMLMessagesSuite) Test_Success(t provider.T) {
	t.Title("Send HTML messages - success")
	t.Severity(allure.NORMAL)
	t.Feature("SMTP sender")
	t.Tags("Positive")
	t.Parallel()

	err := sender.SendHTMLMessages(
		"subject", "<h1>content</h1>", "sender_htmls@mail.ru",
		[]string{"receiver1@mail.ru", "receiver2@mail.ru"},
	)
	t.Require().NoError(err)

	client := resty.New()

	resp, err := client.R().Get(mailhogApiURL + "/api/v2/messages")
	t.Require().NoError(err)
	t.Require().Equal(200, resp.StatusCode())

	messages := GetMessagesResponse{}
	err = json.Unmarshal(resp.Body(), &messages)
	t.Require().NoError(err)

	for _, message := range messages.Items {
		if message.Raw.From == "sender_htmls@mail.ru" {
			t.Require().ElementsMatch([]string{"receiver1@mail.ru", "receiver2@mail.ru"}, message.Raw.To)
			t.Require().Equal("subject", message.Content.Headers.Subject[0])
			t.Require().Equal("<h1>content</h1>", message.Content.Body)
			return
		}
	}

	t.Fail()
}

func (s *SendHTMLMessagesSuite) Test_IncorrectEmail(t provider.T) {
	t.Title("Send HTML messages - incorrect email")
	t.Severity(allure.CRITICAL)
	t.Feature("SMTP sender")
	t.Tags("Negative")
	t.Parallel()

	err := sender.SendHTMLMessages(
		"subject", "content", "sender_html_incorrect_email@mail.ru",
		[]string{"not_email"},
	)
	t.Require().Error(err)
	t.Require().Contains(err.Error(), "invalid address")
}
