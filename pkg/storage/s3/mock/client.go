// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	s3 "github.com/mandarine-io/Backend/pkg/storage/s3"
	mock "github.com/stretchr/testify/mock"
)

// ClientMock is an autogenerated mock type for the Client type
type ClientMock struct {
	mock.Mock
}

type ClientMock_Expecter struct {
	mock *mock.Mock
}

func (_m *ClientMock) EXPECT() *ClientMock_Expecter {
	return &ClientMock_Expecter{mock: &_m.Mock}
}

// CreateMany provides a mock function with given fields: ctx, files
func (_m *ClientMock) CreateMany(ctx context.Context, files []*s3.FileData) map[string]*s3.CreateDto {
	ret := _m.Called(ctx, files)

	if len(ret) == 0 {
		panic("no return value specified for CreateMany")
	}

	var r0 map[string]*s3.CreateDto
	if rf, ok := ret.Get(0).(func(context.Context, []*s3.FileData) map[string]*s3.CreateDto); ok {
		r0 = rf(ctx, files)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*s3.CreateDto)
		}
	}

	return r0
}

// ClientMock_CreateMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMany'
type ClientMock_CreateMany_Call struct {
	*mock.Call
}

// CreateMany is a helper method to define mock.On call
//   - ctx context.Context
//   - files []*s3.FileData
func (_e *ClientMock_Expecter) CreateMany(ctx interface{}, files interface{}) *ClientMock_CreateMany_Call {
	return &ClientMock_CreateMany_Call{Call: _e.mock.On("CreateMany", ctx, files)}
}

func (_c *ClientMock_CreateMany_Call) Run(run func(ctx context.Context, files []*s3.FileData)) *ClientMock_CreateMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*s3.FileData))
	})
	return _c
}

func (_c *ClientMock_CreateMany_Call) Return(_a0 map[string]*s3.CreateDto) *ClientMock_CreateMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_CreateMany_Call) RunAndReturn(run func(context.Context, []*s3.FileData) map[string]*s3.CreateDto) *ClientMock_CreateMany_Call {
	_c.Call.Return(run)
	return _c
}

// CreateOne provides a mock function with given fields: ctx, file
func (_m *ClientMock) CreateOne(ctx context.Context, file *s3.FileData) *s3.CreateDto {
	ret := _m.Called(ctx, file)

	if len(ret) == 0 {
		panic("no return value specified for CreateOne")
	}

	var r0 *s3.CreateDto
	if rf, ok := ret.Get(0).(func(context.Context, *s3.FileData) *s3.CreateDto); ok {
		r0 = rf(ctx, file)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.CreateDto)
		}
	}

	return r0
}

// ClientMock_CreateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateOne'
type ClientMock_CreateOne_Call struct {
	*mock.Call
}

// CreateOne is a helper method to define mock.On call
//   - ctx context.Context
//   - file *s3.FileData
func (_e *ClientMock_Expecter) CreateOne(ctx interface{}, file interface{}) *ClientMock_CreateOne_Call {
	return &ClientMock_CreateOne_Call{Call: _e.mock.On("CreateOne", ctx, file)}
}

func (_c *ClientMock_CreateOne_Call) Run(run func(ctx context.Context, file *s3.FileData)) *ClientMock_CreateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*s3.FileData))
	})
	return _c
}

func (_c *ClientMock_CreateOne_Call) Return(_a0 *s3.CreateDto) *ClientMock_CreateOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_CreateOne_Call) RunAndReturn(run func(context.Context, *s3.FileData) *s3.CreateDto) *ClientMock_CreateOne_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMany provides a mock function with given fields: ctx, objectIDs
func (_m *ClientMock) DeleteMany(ctx context.Context, objectIDs []string) map[string]error {
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

// ClientMock_DeleteMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMany'
type ClientMock_DeleteMany_Call struct {
	*mock.Call
}

// DeleteMany is a helper method to define mock.On call
//   - ctx context.Context
//   - objectIDs []string
func (_e *ClientMock_Expecter) DeleteMany(ctx interface{}, objectIDs interface{}) *ClientMock_DeleteMany_Call {
	return &ClientMock_DeleteMany_Call{Call: _e.mock.On("DeleteMany", ctx, objectIDs)}
}

func (_c *ClientMock_DeleteMany_Call) Run(run func(ctx context.Context, objectIDs []string)) *ClientMock_DeleteMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *ClientMock_DeleteMany_Call) Return(_a0 map[string]error) *ClientMock_DeleteMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_DeleteMany_Call) RunAndReturn(run func(context.Context, []string) map[string]error) *ClientMock_DeleteMany_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOne provides a mock function with given fields: ctx, objectID
func (_m *ClientMock) DeleteOne(ctx context.Context, objectID string) error {
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

// ClientMock_DeleteOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOne'
type ClientMock_DeleteOne_Call struct {
	*mock.Call
}

// DeleteOne is a helper method to define mock.On call
//   - ctx context.Context
//   - objectID string
func (_e *ClientMock_Expecter) DeleteOne(ctx interface{}, objectID interface{}) *ClientMock_DeleteOne_Call {
	return &ClientMock_DeleteOne_Call{Call: _e.mock.On("DeleteOne", ctx, objectID)}
}

func (_c *ClientMock_DeleteOne_Call) Run(run func(ctx context.Context, objectID string)) *ClientMock_DeleteOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ClientMock_DeleteOne_Call) Return(_a0 error) *ClientMock_DeleteOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_DeleteOne_Call) RunAndReturn(run func(context.Context, string) error) *ClientMock_DeleteOne_Call {
	_c.Call.Return(run)
	return _c
}

// GetMany provides a mock function with given fields: ctx, objectIDs
func (_m *ClientMock) GetMany(ctx context.Context, objectIDs []string) map[string]*s3.GetDto {
	ret := _m.Called(ctx, objectIDs)

	if len(ret) == 0 {
		panic("no return value specified for GetMany")
	}

	var r0 map[string]*s3.GetDto
	if rf, ok := ret.Get(0).(func(context.Context, []string) map[string]*s3.GetDto); ok {
		r0 = rf(ctx, objectIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*s3.GetDto)
		}
	}

	return r0
}

// ClientMock_GetMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMany'
type ClientMock_GetMany_Call struct {
	*mock.Call
}

// GetMany is a helper method to define mock.On call
//   - ctx context.Context
//   - objectIDs []string
func (_e *ClientMock_Expecter) GetMany(ctx interface{}, objectIDs interface{}) *ClientMock_GetMany_Call {
	return &ClientMock_GetMany_Call{Call: _e.mock.On("GetMany", ctx, objectIDs)}
}

func (_c *ClientMock_GetMany_Call) Run(run func(ctx context.Context, objectIDs []string)) *ClientMock_GetMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *ClientMock_GetMany_Call) Return(_a0 map[string]*s3.GetDto) *ClientMock_GetMany_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_GetMany_Call) RunAndReturn(run func(context.Context, []string) map[string]*s3.GetDto) *ClientMock_GetMany_Call {
	_c.Call.Return(run)
	return _c
}

// GetOne provides a mock function with given fields: ctx, objectID
func (_m *ClientMock) GetOne(ctx context.Context, objectID string) *s3.GetDto {
	ret := _m.Called(ctx, objectID)

	if len(ret) == 0 {
		panic("no return value specified for GetOne")
	}

	var r0 *s3.GetDto
	if rf, ok := ret.Get(0).(func(context.Context, string) *s3.GetDto); ok {
		r0 = rf(ctx, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.GetDto)
		}
	}

	return r0
}

// ClientMock_GetOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOne'
type ClientMock_GetOne_Call struct {
	*mock.Call
}

// GetOne is a helper method to define mock.On call
//   - ctx context.Context
//   - objectID string
func (_e *ClientMock_Expecter) GetOne(ctx interface{}, objectID interface{}) *ClientMock_GetOne_Call {
	return &ClientMock_GetOne_Call{Call: _e.mock.On("GetOne", ctx, objectID)}
}

func (_c *ClientMock_GetOne_Call) Run(run func(ctx context.Context, objectID string)) *ClientMock_GetOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *ClientMock_GetOne_Call) Return(_a0 *s3.GetDto) *ClientMock_GetOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ClientMock_GetOne_Call) RunAndReturn(run func(context.Context, string) *s3.GetDto) *ClientMock_GetOne_Call {
	_c.Call.Return(run)
	return _c
}

// NewClientMock creates a new instance of ClientMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientMock {
	mock := &ClientMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
