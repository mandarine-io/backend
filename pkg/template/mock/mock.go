package mock

import (
	"github.com/stretchr/testify/mock"
)

type TemplateEngineMock struct {
	mock.Mock
}

func (m *TemplateEngineMock) Render(templateName string, templateArgs any) (string, error) {
	args := m.Called(templateName, templateArgs)
	return args.String(0), args.Error(1)
}
