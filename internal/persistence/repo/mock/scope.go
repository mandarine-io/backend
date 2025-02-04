// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"
	gorm "gorm.io/gorm"
)

// ScopeMock is an autogenerated mock type for the Scope type
type ScopeMock struct {
	mock.Mock
}

type ScopeMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ScopeMock) EXPECT() *ScopeMock_Expecter {
	return &ScopeMock_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: db
func (_m *ScopeMock) Execute(db *gorm.DB) *gorm.DB {
	ret := _m.Called(db)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB) *gorm.DB); ok {
		r0 = rf(db)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// ScopeMock_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type ScopeMock_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - db *gorm.DB
func (_e *ScopeMock_Expecter) Execute(db interface{}) *ScopeMock_Execute_Call {
	return &ScopeMock_Execute_Call{Call: _e.mock.On("Execute", db)}
}

func (_c *ScopeMock_Execute_Call) Run(run func(db *gorm.DB)) *ScopeMock_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gorm.DB))
	})
	return _c
}

func (_c *ScopeMock_Execute_Call) Return(_a0 *gorm.DB) *ScopeMock_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ScopeMock_Execute_Call) RunAndReturn(run func(*gorm.DB) *gorm.DB) *ScopeMock_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewScopeMock creates a new instance of ScopeMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewScopeMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ScopeMock {
	mock := &ScopeMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
