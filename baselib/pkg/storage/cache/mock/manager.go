// Code generated by mockery v2.50.1. DO NOT EDIT.

package mock

import (
	context "context"
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// ManagerMock is an autogenerated mock type for the Manager type
type ManagerMock struct {
	mock.Mock
}

type ManagerMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ManagerMock) EXPECT() *ManagerMock_Expecter {
	return &ManagerMock_Expecter{mock: &_m.Mock}
}

// Delete provides a mock function with given fields: ctx, keys
func (_m *ManagerMock) Delete(ctx context.Context, keys ...string) error {
	_va := make([]interface{}, len(keys))
	for _i := range keys {
		_va[_i] = keys[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) error); ok {
		r0 = rf(ctx, keys...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type ManagerMock_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - keys ...string
func (_e *ManagerMock_Expecter) Delete(ctx interface{}, keys ...interface{}) *ManagerMock_Delete_Call {
	return &ManagerMock_Delete_Call{Call: _e.mock.On("Delete",
		append([]interface{}{ctx}, keys...)...)}
}

func (_c *ManagerMock_Delete_Call) Run(run func(ctx context.Context, keys ...string)) *ManagerMock_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *ManagerMock_Delete_Call) Return(_a0 error) *ManagerMock_Delete_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_Delete_Call) RunAndReturn(run func(context.Context, ...string) error) *ManagerMock_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, key, value
func (_m *ManagerMock) Get(ctx context.Context, key string, value interface{}) error {
	ret := _m.Called(ctx, key, value)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ManagerMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
func (_e *ManagerMock_Expecter) Get(ctx interface{}, key interface{}, value interface{}) *ManagerMock_Get_Call {
	return &ManagerMock_Get_Call{Call: _e.mock.On("Get", ctx, key, value)}
}

func (_c *ManagerMock_Get_Call) Run(run func(ctx context.Context, key string, value interface{})) *ManagerMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *ManagerMock_Get_Call) Return(_a0 error) *ManagerMock_Get_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_Get_Call) RunAndReturn(run func(context.Context, string, interface{}) error) *ManagerMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Invalidate provides a mock function with given fields: ctx, keyRegex
func (_m *ManagerMock) Invalidate(ctx context.Context, keyRegex string) error {
	ret := _m.Called(ctx, keyRegex)

	if len(ret) == 0 {
		panic("no return value specified for Invalidate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, keyRegex)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_Invalidate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Invalidate'
type ManagerMock_Invalidate_Call struct {
	*mock.Call
}

// Invalidate is a helper method to define mock.On call
//   - ctx context.Context
//   - keyRegex string
func (_e *ManagerMock_Expecter) Invalidate(ctx interface{}, keyRegex interface{}) *ManagerMock_Invalidate_Call {
	return &ManagerMock_Invalidate_Call{Call: _e.mock.On("Invalidate", ctx, keyRegex)}
}

func (_c *ManagerMock_Invalidate_Call) Run(run func(ctx context.Context, keyRegex string)) *ManagerMock_Invalidate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ManagerMock_Invalidate_Call) Return(_a0 error) *ManagerMock_Invalidate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_Invalidate_Call) RunAndReturn(run func(context.Context, string) error) *ManagerMock_Invalidate_Call {
	_c.Call.Return(run)
	return _c
}

// Set provides a mock function with given fields: ctx, key, value
func (_m *ManagerMock) Set(ctx context.Context, key string, value interface{}) error {
	ret := _m.Called(ctx, key, value)

	if len(ret) == 0 {
		panic("no return value specified for Set")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_Set_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Set'
type ManagerMock_Set_Call struct {
	*mock.Call
}

// Set is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
func (_e *ManagerMock_Expecter) Set(ctx interface{}, key interface{}, value interface{}) *ManagerMock_Set_Call {
	return &ManagerMock_Set_Call{Call: _e.mock.On("Set", ctx, key, value)}
}

func (_c *ManagerMock_Set_Call) Run(run func(ctx context.Context, key string, value interface{})) *ManagerMock_Set_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *ManagerMock_Set_Call) Return(_a0 error) *ManagerMock_Set_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_Set_Call) RunAndReturn(run func(context.Context, string, interface{}) error) *ManagerMock_Set_Call {
	_c.Call.Return(run)
	return _c
}

// SetWithExpiration provides a mock function with given fields: ctx, key, value, expiration
func (_m *ManagerMock) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ret := _m.Called(ctx, key, value, expiration)

	if len(ret) == 0 {
		panic("no return value specified for SetWithExpiration")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_SetWithExpiration_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SetWithExpiration'
type ManagerMock_SetWithExpiration_Call struct {
	*mock.Call
}

// SetWithExpiration is a helper method to define mock.On call
//   - ctx context.Context
//   - key string
//   - value interface{}
//   - expiration time.Duration
func (_e *ManagerMock_Expecter) SetWithExpiration(ctx interface{}, key interface{}, value interface{}, expiration interface{}) *ManagerMock_SetWithExpiration_Call {
	return &ManagerMock_SetWithExpiration_Call{Call: _e.mock.On("SetWithExpiration", ctx, key, value, expiration)}
}

func (_c *ManagerMock_SetWithExpiration_Call) Run(run func(ctx context.Context, key string, value interface{}, expiration time.Duration)) *ManagerMock_SetWithExpiration_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}), args[3].(time.Duration))
	})
	return _c
}

func (_c *ManagerMock_SetWithExpiration_Call) Return(_a0 error) *ManagerMock_SetWithExpiration_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_SetWithExpiration_Call) RunAndReturn(run func(context.Context, string, interface{}, time.Duration) error) *ManagerMock_SetWithExpiration_Call {
	_c.Call.Return(run)
	return _c
}

// NewManagerMock creates a new instance of ManagerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewManagerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ManagerMock {
	mock := &ManagerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}