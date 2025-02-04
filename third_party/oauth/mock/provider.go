// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	oauth "github.com/mandarine-io/backend/third_party/oauth"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// ProviderMock is an autogenerated mock type for the Provider type
type ProviderMock struct {
	mock.Mock
}

type ProviderMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ProviderMock) EXPECT() *ProviderMock_Expecter {
	return &ProviderMock_Expecter{mock: &_m.Mock}
}

// ExchangeCodeToToken provides a mock function with given fields: ctx, code, redirectURL
func (_m *ProviderMock) ExchangeCodeToToken(ctx context.Context, code string, redirectURL string) (*oauth2.Token, error) {
	ret := _m.Called(ctx, code, redirectURL)

	if len(ret) == 0 {
		panic("no return value specified for ExchangeCodeToToken")
	}

	var r0 *oauth2.Token
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*oauth2.Token, error)); ok {
		return rf(ctx, code, redirectURL)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *oauth2.Token); ok {
		r0 = rf(ctx, code, redirectURL)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, code, redirectURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderMock_ExchangeCodeToToken_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExchangeCodeToToken'
type ProviderMock_ExchangeCodeToToken_Call struct {
	*mock.Call
}

// ExchangeCodeToToken is a helper method to define mock.On call
//   - ctx context.Context
//   - code string
//   - redirectURL string
func (_e *ProviderMock_Expecter) ExchangeCodeToToken(ctx interface{}, code interface{}, redirectURL interface{}) *ProviderMock_ExchangeCodeToToken_Call {
	return &ProviderMock_ExchangeCodeToToken_Call{Call: _e.mock.On("ExchangeCodeToToken", ctx, code, redirectURL)}
}

func (_c *ProviderMock_ExchangeCodeToToken_Call) Run(run func(ctx context.Context, code string, redirectURL string)) *ProviderMock_ExchangeCodeToToken_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *ProviderMock_ExchangeCodeToToken_Call) Return(_a0 *oauth2.Token, _a1 error) *ProviderMock_ExchangeCodeToToken_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderMock_ExchangeCodeToToken_Call) RunAndReturn(run func(context.Context, string, string) (*oauth2.Token, error)) *ProviderMock_ExchangeCodeToToken_Call {
	_c.Call.Return(run)
	return _c
}

// GetConsentPageURL provides a mock function with given fields: redirectURL
func (_m *ProviderMock) GetConsentPageURL(redirectURL string) (string, string) {
	ret := _m.Called(redirectURL)

	if len(ret) == 0 {
		panic("no return value specified for GetConsentPageURL")
	}

	var r0 string
	var r1 string
	if rf, ok := ret.Get(0).(func(string) (string, string)); ok {
		return rf(redirectURL)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(redirectURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(redirectURL)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ProviderMock_GetConsentPageURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetConsentPageURL'
type ProviderMock_GetConsentPageURL_Call struct {
	*mock.Call
}

// GetConsentPageURL is a helper method to define mock.On call
//   - redirectURL string
func (_e *ProviderMock_Expecter) GetConsentPageURL(redirectURL interface{}) *ProviderMock_GetConsentPageURL_Call {
	return &ProviderMock_GetConsentPageURL_Call{Call: _e.mock.On("GetConsentPageURL", redirectURL)}
}

func (_c *ProviderMock_GetConsentPageURL_Call) Run(run func(redirectURL string)) *ProviderMock_GetConsentPageURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ProviderMock_GetConsentPageURL_Call) Return(_a0 string, _a1 string) *ProviderMock_GetConsentPageURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderMock_GetConsentPageURL_Call) RunAndReturn(run func(string) (string, string)) *ProviderMock_GetConsentPageURL_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserInfo provides a mock function with given fields: ctx, token
func (_m *ProviderMock) GetUserInfo(ctx context.Context, token *oauth2.Token) (oauth.UserInfo, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for GetUserInfo")
	}

	var r0 oauth.UserInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *oauth2.Token) (oauth.UserInfo, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *oauth2.Token) oauth.UserInfo); ok {
		r0 = rf(ctx, token)
	} else {
		r0 = ret.Get(0).(oauth.UserInfo)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *oauth2.Token) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderMock_GetUserInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserInfo'
type ProviderMock_GetUserInfo_Call struct {
	*mock.Call
}

// GetUserInfo is a helper method to define mock.On call
//   - ctx context.Context
//   - token *oauth2.Token
func (_e *ProviderMock_Expecter) GetUserInfo(ctx interface{}, token interface{}) *ProviderMock_GetUserInfo_Call {
	return &ProviderMock_GetUserInfo_Call{Call: _e.mock.On("GetUserInfo", ctx, token)}
}

func (_c *ProviderMock_GetUserInfo_Call) Run(run func(ctx context.Context, token *oauth2.Token)) *ProviderMock_GetUserInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*oauth2.Token))
	})
	return _c
}

func (_c *ProviderMock_GetUserInfo_Call) Return(_a0 oauth.UserInfo, _a1 error) *ProviderMock_GetUserInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderMock_GetUserInfo_Call) RunAndReturn(run func(context.Context, *oauth2.Token) (oauth.UserInfo, error)) *ProviderMock_GetUserInfo_Call {
	_c.Call.Return(run)
	return _c
}

// NewProviderMock creates a new instance of ProviderMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProviderMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProviderMock {
	mock := &ProviderMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
