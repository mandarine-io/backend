package mock

import "github.com/stretchr/testify/mock"

type SenderMock struct {
	mock.Mock
}

func (m *SenderMock) HealthCheck() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *SenderMock) SendPlainMessage(subject string, content string, to string, attachments ...string) error {
	args := m.Called(subject, content, to, attachments)
	return args.Error(0)
}

func (m *SenderMock) SendPlainMessages(subject string, content string, to []string, attachments ...string) error {
	args := m.Called(subject, content, to, attachments)
	return args.Error(0)
}

func (m *SenderMock) SendHtmlMessage(subject string, content string, to string, attachments ...string) error {
	args := m.Called(subject, content, to, attachments)
	return args.Error(0)
}

func (m *SenderMock) SendHtmlMessages(subject string, content string, to []string, attachments ...string) error {
	args := m.Called(subject, content, to, attachments)
	return args.Error(0)
}
