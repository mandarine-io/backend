// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	s3 "github.com/mandarine-io/backend/internal/infrastructure/s3"
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

// CreateMany provides a mock function with given fields: ctx, files
func (_m *ManagerMock) CreateMany(ctx context.Context, files []*s3.FileData) map[string]s3.CreateResult {
	ret := _m.Called(ctx, files)

	if len(ret) == 0 {
		panic("no return value specified for CreateMany")
	}

	var r0 map[string]s3.CreateResult
	if rf, ok := ret.Get(0).(func(context.Context, []*s3.FileData) map[string]s3.CreateResult); ok {
		r0 = rf(ctx, files)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]s3.CreateResult)
		}
	}

	return r0
}

// ManagerMock_CreateMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMany'
type ManagerMock_CreateMany_Call struct {
	*mock.Call
}

// CreateMany is a helper method to define mock.On call
//   - ctx context.Context
//   - files []*s3.FileData
func (_e *ManagerMock_Expecter) CreateMany(ctx interface{}, files interface{}) *ManagerMock_CreateMany_Call {
	return &ManagerMock_CreateMany_Call{Call: _e.mock.On("CreateMany", ctx, files)}
}

func (_c *ManagerMock_CreateMany_Call) Run(run func(ctx context.Context, files []*s3.FileData)) *ManagerMock_CreateMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*s3.FileData))
	})
	return _c
}

func (_c *ManagerMock_CreateMany_Call) Return(_a0 map[string]s3.CreateResult) *ManagerMock_CreateMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_CreateMany_Call) RunAndReturn(run func(context.Context, []*s3.FileData) map[string]s3.CreateResult) *ManagerMock_CreateMany_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOne provides a mock function with given fields: ctx, file
func (_m *ManagerMock) CreateOne(ctx context.Context, file *s3.FileData) s3.CreateResult {
	ret := _m.Called(ctx, file)

	if len(ret) == 0 {
		panic("no return value specified for CreateOne")
	}

	var r0 s3.CreateResult
	if rf, ok := ret.Get(0).(func(context.Context, *s3.FileData) s3.CreateResult); ok {
		r0 = rf(ctx, file)
	} else {
		r0 = ret.Get(0).(s3.CreateResult)
	}

	return r0
}

// ManagerMock_CreateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOne'
type ManagerMock_CreateOne_Call struct {
	*mock.Call
}

// CreateOne is a helper method to define mock.On call
//   - ctx context.Context
//   - file *s3.FileData
func (_e *ManagerMock_Expecter) CreateOne(ctx interface{}, file interface{}) *ManagerMock_CreateOne_Call {
	return &ManagerMock_CreateOne_Call{Call: _e.mock.On("CreateOne", ctx, file)}
}

func (_c *ManagerMock_CreateOne_Call) Run(run func(ctx context.Context, file *s3.FileData)) *ManagerMock_CreateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*s3.FileData))
	})
	return _c
}

func (_c *ManagerMock_CreateOne_Call) Return(_a0 s3.CreateResult) *ManagerMock_CreateOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_CreateOne_Call) RunAndReturn(run func(context.Context, *s3.FileData) s3.CreateResult) *ManagerMock_CreateOne_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMany provides a mock function with given fields: ctx, objectIDs
func (_m *ManagerMock) DeleteMany(ctx context.Context, objectIDs []string) map[string]error {
	ret := _m.Called(ctx, objectIDs)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMany")
	}

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]error); ok {
		r0 = rf(ctx, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// ManagerMock_DeleteMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMany'
type ManagerMock_DeleteMany_Call struct {
	*mock.Call
}

// DeleteMany is a helper method to define mock.On call
//   - ctx context.Context
//   - objectIDs []string
func (_e *ManagerMock_Expecter) DeleteMany(ctx interface{}, objectIDs interface{}) *ManagerMock_DeleteMany_Call {
	return &ManagerMock_DeleteMany_Call{Call: _e.mock.On("DeleteMany", ctx, objectIDs)}
}

func (_c *ManagerMock_DeleteMany_Call) Run(run func(ctx context.Context, objectIDs []string)) *ManagerMock_DeleteMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *ManagerMock_DeleteMany_Call) Return(_a0 map[string]error) *ManagerMock_DeleteMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_DeleteMany_Call) RunAndReturn(run func(context.Context, []string) map[string]error) *ManagerMock_DeleteMany_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOne provides a mock function with given fields: ctx, objectID
func (_m *ManagerMock) DeleteOne(ctx context.Context, objectID string) error {
	ret := _m.Called(ctx, objectID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, objectID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ManagerMock_DeleteOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOne'
type ManagerMock_DeleteOne_Call struct {
	*mock.Call
}

// DeleteOne is a helper method to define mock.On call
//   - ctx context.Context
//   - objectID string
func (_e *ManagerMock_Expecter) DeleteOne(ctx interface{}, objectID interface{}) *ManagerMock_DeleteOne_Call {
	return &ManagerMock_DeleteOne_Call{Call: _e.mock.On("DeleteOne", ctx, objectID)}
}

func (_c *ManagerMock_DeleteOne_Call) Run(run func(ctx context.Context, objectID string)) *ManagerMock_DeleteOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ManagerMock_DeleteOne_Call) Return(_a0 error) *ManagerMock_DeleteOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_DeleteOne_Call) RunAndReturn(run func(context.Context, string) error) *ManagerMock_DeleteOne_Call {
	_c.Call.Return(run)
	return _c
}

// GetMany provides a mock function with given fields: ctx, objectIDs
func (_m *ManagerMock) GetMany(ctx context.Context, objectIDs []string) map[string]s3.GetResult {
	ret := _m.Called(ctx, objectIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetMany")
	}

	var r0 map[string]s3.GetResult
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]s3.GetResult); ok {
		r0 = rf(ctx, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]s3.GetResult)
		}
	}

	return r0
}

// ManagerMock_GetMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMany'
type ManagerMock_GetMany_Call struct {
	*mock.Call
}

// GetMany is a helper method to define mock.On call
//   - ctx context.Context
//   - objectIDs []string
func (_e *ManagerMock_Expecter) GetMany(ctx interface{}, objectIDs interface{}) *ManagerMock_GetMany_Call {
	return &ManagerMock_GetMany_Call{Call: _e.mock.On("GetMany", ctx, objectIDs)}
}

func (_c *ManagerMock_GetMany_Call) Run(run func(ctx context.Context, objectIDs []string)) *ManagerMock_GetMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *ManagerMock_GetMany_Call) Return(_a0 map[string]s3.GetResult) *ManagerMock_GetMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_GetMany_Call) RunAndReturn(run func(context.Context, []string) map[string]s3.GetResult) *ManagerMock_GetMany_Call {
	_c.Call.Return(run)
	return _c
}

// GetOne provides a mock function with given fields: ctx, objectID
func (_m *ManagerMock) GetOne(ctx context.Context, objectID string) s3.GetResult {
	ret := _m.Called(ctx, objectID)

	if len(ret) == 0 {
		panic("no return value specified for GetOne")
	}

	var r0 s3.GetResult
	if rf, ok := ret.Get(0).(func(context.Context, string) s3.GetResult); ok {
		r0 = rf(ctx, objectID)
	} else {
		r0 = ret.Get(0).(s3.GetResult)
	}

	return r0
}

// ManagerMock_GetOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOne'
type ManagerMock_GetOne_Call struct {
	*mock.Call
}

// GetOne is a helper method to define mock.On call
//   - ctx context.Context
//   - objectID string
func (_e *ManagerMock_Expecter) GetOne(ctx interface{}, objectID interface{}) *ManagerMock_GetOne_Call {
	return &ManagerMock_GetOne_Call{Call: _e.mock.On("GetOne", ctx, objectID)}
}

func (_c *ManagerMock_GetOne_Call) Run(run func(ctx context.Context, objectID string)) *ManagerMock_GetOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ManagerMock_GetOne_Call) Return(_a0 s3.GetResult) *ManagerMock_GetOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ManagerMock_GetOne_Call) RunAndReturn(run func(context.Context, string) s3.GetResult) *ManagerMock_GetOne_Call {
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
