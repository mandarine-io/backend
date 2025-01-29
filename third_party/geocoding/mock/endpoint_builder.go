// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	geocoding "github.com/mandarine-io/backend/third_party/geocoding"
	mock "github.com/stretchr/testify/mock"
)

// EndpointBuilderMock is an autogenerated mock type for the EndpointBuilder type
type EndpointBuilderMock struct {
	mock.Mock
}

type EndpointBuilderMock_Expecter struct {
	mock *mock.Mock
}

func (_m *EndpointBuilderMock) EXPECT() *EndpointBuilderMock_Expecter {
	return &EndpointBuilderMock_Expecter{mock: &_m.Mock}
}

// GeocodeURL provides a mock function with given fields: address, cfg
func (_m *EndpointBuilderMock) GeocodeURL(address string, cfg geocoding.GeocodeConfig) string {
	ret := _m.Called(address, cfg)

	if len(ret) == 0 {
		panic("no return value specified for GeocodeURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(string, geocoding.GeocodeConfig) string); ok {
		r0 = rf(address, cfg)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EndpointBuilderMock_GeocodeURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GeocodeURL'
type EndpointBuilderMock_GeocodeURL_Call struct {
	*mock.Call
}

// GeocodeURL is a helper method to define mock.On call
//   - address string
//   - cfg geocoding.GeocodeConfig
func (_e *EndpointBuilderMock_Expecter) GeocodeURL(address interface{}, cfg interface{}) *EndpointBuilderMock_GeocodeURL_Call {
	return &EndpointBuilderMock_GeocodeURL_Call{Call: _e.mock.On("GeocodeURL", address, cfg)}
}

func (_c *EndpointBuilderMock_GeocodeURL_Call) Run(run func(address string, cfg geocoding.GeocodeConfig)) *EndpointBuilderMock_GeocodeURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(geocoding.GeocodeConfig))
	})
	return _c
}

func (_c *EndpointBuilderMock_GeocodeURL_Call) Return(_a0 string) *EndpointBuilderMock_GeocodeURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EndpointBuilderMock_GeocodeURL_Call) RunAndReturn(run func(string, geocoding.GeocodeConfig) string) *EndpointBuilderMock_GeocodeURL_Call {
	_c.Call.Return(run)
	return _c
}

// ReverseGeocodeURL provides a mock function with given fields: loc, cfg
func (_m *EndpointBuilderMock) ReverseGeocodeURL(loc geocoding.Location, cfg geocoding.ReverseGeocodeConfig) string {
	ret := _m.Called(loc, cfg)

	if len(ret) == 0 {
		panic("no return value specified for ReverseGeocodeURL")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func(geocoding.Location, geocoding.ReverseGeocodeConfig) string); ok {
		r0 = rf(loc, cfg)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// EndpointBuilderMock_ReverseGeocodeURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReverseGeocodeURL'
type EndpointBuilderMock_ReverseGeocodeURL_Call struct {
	*mock.Call
}

// ReverseGeocodeURL is a helper method to define mock.On call
//   - loc geocoding.Location
//   - cfg geocoding.ReverseGeocodeConfig
func (_e *EndpointBuilderMock_Expecter) ReverseGeocodeURL(loc interface{}, cfg interface{}) *EndpointBuilderMock_ReverseGeocodeURL_Call {
	return &EndpointBuilderMock_ReverseGeocodeURL_Call{Call: _e.mock.On("ReverseGeocodeURL", loc, cfg)}
}

func (_c *EndpointBuilderMock_ReverseGeocodeURL_Call) Run(run func(loc geocoding.Location, cfg geocoding.ReverseGeocodeConfig)) *EndpointBuilderMock_ReverseGeocodeURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(geocoding.Location), args[1].(geocoding.ReverseGeocodeConfig))
	})
	return _c
}

func (_c *EndpointBuilderMock_ReverseGeocodeURL_Call) Return(_a0 string) *EndpointBuilderMock_ReverseGeocodeURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EndpointBuilderMock_ReverseGeocodeURL_Call) RunAndReturn(run func(geocoding.Location, geocoding.ReverseGeocodeConfig) string) *EndpointBuilderMock_ReverseGeocodeURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewEndpointBuilderMock creates a new instance of EndpointBuilderMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEndpointBuilderMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *EndpointBuilderMock {
	mock := &EndpointBuilderMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
