// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	geocoding "github.com/mandarine-io/backend/third_party/geocoding"
	mock "github.com/stretchr/testify/mock"
)

// ResponseParserMock is an autogenerated mock type for the ResponseParser type
type ResponseParserMock struct {
	mock.Mock
}

type ResponseParserMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ResponseParserMock) EXPECT() *ResponseParserMock_Expecter {
	return &ResponseParserMock_Expecter{mock: &_m.Mock}
}

// Addresses provides a mock function with no fields
func (_m *ResponseParserMock) Addresses() ([]*geocoding.Address, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Addresses")
	}

	var r0 []*geocoding.Address
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*geocoding.Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*geocoding.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*geocoding.Address)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResponseParserMock_Addresses_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Addresses'
type ResponseParserMock_Addresses_Call struct {
	*mock.Call
}

// Addresses is a helper method to define mock.On call
func (_e *ResponseParserMock_Expecter) Addresses() *ResponseParserMock_Addresses_Call {
	return &ResponseParserMock_Addresses_Call{Call: _e.mock.On("Addresses")}
}

func (_c *ResponseParserMock_Addresses_Call) Run(run func()) *ResponseParserMock_Addresses_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ResponseParserMock_Addresses_Call) Return(_a0 []*geocoding.Address, _a1 error) *ResponseParserMock_Addresses_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResponseParserMock_Addresses_Call) RunAndReturn(run func() ([]*geocoding.Address, error)) *ResponseParserMock_Addresses_Call {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with no fields
func (_m *ResponseParserMock) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResponseParserMock_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type ResponseParserMock_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *ResponseParserMock_Expecter) Error() *ResponseParserMock_Error_Call {
	return &ResponseParserMock_Error_Call{Call: _e.mock.On("Error")}
}

func (_c *ResponseParserMock_Error_Call) Run(run func()) *ResponseParserMock_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ResponseParserMock_Error_Call) Return(_a0 error) *ResponseParserMock_Error_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ResponseParserMock_Error_Call) RunAndReturn(run func() error) *ResponseParserMock_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Locations provides a mock function with no fields
func (_m *ResponseParserMock) Locations() ([]*geocoding.Location, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Locations")
	}

	var r0 []*geocoding.Location
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]*geocoding.Location, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []*geocoding.Location); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*geocoding.Location)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResponseParserMock_Locations_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Locations'
type ResponseParserMock_Locations_Call struct {
	*mock.Call
}

// Locations is a helper method to define mock.On call
func (_e *ResponseParserMock_Expecter) Locations() *ResponseParserMock_Locations_Call {
	return &ResponseParserMock_Locations_Call{Call: _e.mock.On("Locations")}
}

func (_c *ResponseParserMock_Locations_Call) Run(run func()) *ResponseParserMock_Locations_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ResponseParserMock_Locations_Call) Return(_a0 []*geocoding.Location, _a1 error) *ResponseParserMock_Locations_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResponseParserMock_Locations_Call) RunAndReturn(run func() ([]*geocoding.Location, error)) *ResponseParserMock_Locations_Call {
	_c.Call.Return(run)
	return _c
}

// NewResponseParserMock creates a new instance of ResponseParserMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResponseParserMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResponseParserMock {
	mock := &ResponseParserMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
