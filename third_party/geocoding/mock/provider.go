// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	geocoding "github.com/mandarine-io/backend/third_party/geocoding"
	mock "github.com/stretchr/testify/mock"
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

// Geocode provides a mock function with given fields: ctx, address, config
func (_m *ProviderMock) Geocode(ctx context.Context, address string, config geocoding.GeocodeConfig) ([]*geocoding.Location, error) {
	ret := _m.Called(ctx, address, config)

	if len(ret) == 0 {
		panic("no return value specified for Geocode")
	}

	var r0 []*geocoding.Location
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, geocoding.GeocodeConfig) ([]*geocoding.Location, error)); ok {
		return rf(ctx, address, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, geocoding.GeocodeConfig) []*geocoding.Location); ok {
		r0 = rf(ctx, address, config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*geocoding.Location)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, geocoding.GeocodeConfig) error); ok {
		r1 = rf(ctx, address, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderMock_Geocode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Geocode'
type ProviderMock_Geocode_Call struct {
	*mock.Call
}

// Geocode is a helper method to define mock.On call
//   - ctx context.Context
//   - address string
//   - config geocoding.GeocodeConfig
func (_e *ProviderMock_Expecter) Geocode(ctx interface{}, address interface{}, config interface{}) *ProviderMock_Geocode_Call {
	return &ProviderMock_Geocode_Call{Call: _e.mock.On("Geocode", ctx, address, config)}
}

func (_c *ProviderMock_Geocode_Call) Run(run func(ctx context.Context, address string, config geocoding.GeocodeConfig)) *ProviderMock_Geocode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(geocoding.GeocodeConfig))
	})
	return _c
}

func (_c *ProviderMock_Geocode_Call) Return(_a0 []*geocoding.Location, _a1 error) *ProviderMock_Geocode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderMock_Geocode_Call) RunAndReturn(run func(context.Context, string, geocoding.GeocodeConfig) ([]*geocoding.Location, error)) *ProviderMock_Geocode_Call {
	_c.Call.Return(run)
	return _c
}

// ReverseGeocode provides a mock function with given fields: ctx, loc, config
func (_m *ProviderMock) ReverseGeocode(ctx context.Context, loc geocoding.Location, config geocoding.ReverseGeocodeConfig) ([]*geocoding.Address, error) {
	ret := _m.Called(ctx, loc, config)

	if len(ret) == 0 {
		panic("no return value specified for ReverseGeocode")
	}

	var r0 []*geocoding.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, geocoding.Location, geocoding.ReverseGeocodeConfig) ([]*geocoding.Address, error)); ok {
		return rf(ctx, loc, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, geocoding.Location, geocoding.ReverseGeocodeConfig) []*geocoding.Address); ok {
		r0 = rf(ctx, loc, config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*geocoding.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, geocoding.Location, geocoding.ReverseGeocodeConfig) error); ok {
		r1 = rf(ctx, loc, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProviderMock_ReverseGeocode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReverseGeocode'
type ProviderMock_ReverseGeocode_Call struct {
	*mock.Call
}

// ReverseGeocode is a helper method to define mock.On call
//   - ctx context.Context
//   - loc geocoding.Location
//   - config geocoding.ReverseGeocodeConfig
func (_e *ProviderMock_Expecter) ReverseGeocode(ctx interface{}, loc interface{}, config interface{}) *ProviderMock_ReverseGeocode_Call {
	return &ProviderMock_ReverseGeocode_Call{Call: _e.mock.On("ReverseGeocode", ctx, loc, config)}
}

func (_c *ProviderMock_ReverseGeocode_Call) Run(run func(ctx context.Context, loc geocoding.Location, config geocoding.ReverseGeocodeConfig)) *ProviderMock_ReverseGeocode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(geocoding.Location), args[2].(geocoding.ReverseGeocodeConfig))
	})
	return _c
}

func (_c *ProviderMock_ReverseGeocode_Call) Return(_a0 []*geocoding.Address, _a1 error) *ProviderMock_ReverseGeocode_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ProviderMock_ReverseGeocode_Call) RunAndReturn(run func(context.Context, geocoding.Location, geocoding.ReverseGeocodeConfig) ([]*geocoding.Address, error)) *ProviderMock_ReverseGeocode_Call {
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
