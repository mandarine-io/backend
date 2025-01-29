// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"
	model "github.com/mandarine-io/backend/pkg/model/v0"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MasterServiceServiceMock is an autogenerated mock type for the MasterServiceService type
type MasterServiceServiceMock struct {
	mock.Mock
}

type MasterServiceServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MasterServiceServiceMock) EXPECT() *MasterServiceServiceMock_Expecter {
	return &MasterServiceServiceMock_Expecter{mock: &_m.Mock}
}

// CreateMasterService provides a mock function with given fields: ctx, userID, input
func (_m *MasterServiceServiceMock) CreateMasterService(ctx context.Context, userID uuid.UUID, input model.CreateMasterServiceInput) (model.MasterServiceOutput, error) {
	ret := _m.Called(ctx, userID, input)

	if len(ret) == 0 {
		panic("no return value specified for CreateMasterService")
	}

	var r0 model.MasterServiceOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.CreateMasterServiceInput) (model.MasterServiceOutput, error)); ok {
		return rf(ctx, userID, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.CreateMasterServiceInput) model.MasterServiceOutput); ok {
		r0 = rf(ctx, userID, input)
	} else {
		r0 = ret.Get(0).(model.MasterServiceOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, model.CreateMasterServiceInput) error); ok {
		r1 = rf(ctx, userID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_CreateMasterService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMasterService'
type MasterServiceServiceMock_CreateMasterService_Call struct {
	*mock.Call
}

// CreateMasterService is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - input model.CreateMasterServiceInput
func (_e *MasterServiceServiceMock_Expecter) CreateMasterService(ctx interface{}, userID interface{}, input interface{}) *MasterServiceServiceMock_CreateMasterService_Call {
	return &MasterServiceServiceMock_CreateMasterService_Call{Call: _e.mock.On("CreateMasterService", ctx, userID, input)}
}

func (_c *MasterServiceServiceMock_CreateMasterService_Call) Run(run func(ctx context.Context, userID uuid.UUID, input model.CreateMasterServiceInput)) *MasterServiceServiceMock_CreateMasterService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(model.CreateMasterServiceInput))
	})
	return _c
}

func (_c *MasterServiceServiceMock_CreateMasterService_Call) Return(_a0 model.MasterServiceOutput, _a1 error) *MasterServiceServiceMock_CreateMasterService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_CreateMasterService_Call) RunAndReturn(run func(context.Context, uuid.UUID, model.CreateMasterServiceInput) (model.MasterServiceOutput, error)) *MasterServiceServiceMock_CreateMasterService_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMasterService provides a mock function with given fields: ctx, userID, masterServiceID
func (_m *MasterServiceServiceMock) DeleteMasterService(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID) error {
	ret := _m.Called(ctx, userID, masterServiceID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMasterService")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, userID, masterServiceID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MasterServiceServiceMock_DeleteMasterService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMasterService'
type MasterServiceServiceMock_DeleteMasterService_Call struct {
	*mock.Call
}

// DeleteMasterService is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - masterServiceID uuid.UUID
func (_e *MasterServiceServiceMock_Expecter) DeleteMasterService(ctx interface{}, userID interface{}, masterServiceID interface{}) *MasterServiceServiceMock_DeleteMasterService_Call {
	return &MasterServiceServiceMock_DeleteMasterService_Call{Call: _e.mock.On("DeleteMasterService", ctx, userID, masterServiceID)}
}

func (_c *MasterServiceServiceMock_DeleteMasterService_Call) Run(run func(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID)) *MasterServiceServiceMock_DeleteMasterService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceServiceMock_DeleteMasterService_Call) Return(_a0 error) *MasterServiceServiceMock_DeleteMasterService_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceServiceMock_DeleteMasterService_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MasterServiceServiceMock_DeleteMasterService_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllMasterServices provides a mock function with given fields: ctx, input
func (_m *MasterServiceServiceMock) FindAllMasterServices(ctx context.Context, input model.FindMasterServicesInput) (model.MasterServicesOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for FindAllMasterServices")
	}

	var r0 model.MasterServicesOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, model.FindMasterServicesInput) (model.MasterServicesOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, model.FindMasterServicesInput) model.MasterServicesOutput); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(model.MasterServicesOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, model.FindMasterServicesInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_FindAllMasterServices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllMasterServices'
type MasterServiceServiceMock_FindAllMasterServices_Call struct {
	*mock.Call
}

// FindAllMasterServices is a helper method to define mock.On call
//   - ctx context.Context
//   - input model.FindMasterServicesInput
func (_e *MasterServiceServiceMock_Expecter) FindAllMasterServices(ctx interface{}, input interface{}) *MasterServiceServiceMock_FindAllMasterServices_Call {
	return &MasterServiceServiceMock_FindAllMasterServices_Call{Call: _e.mock.On("FindAllMasterServices", ctx, input)}
}

func (_c *MasterServiceServiceMock_FindAllMasterServices_Call) Run(run func(ctx context.Context, input model.FindMasterServicesInput)) *MasterServiceServiceMock_FindAllMasterServices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.FindMasterServicesInput))
	})
	return _c
}

func (_c *MasterServiceServiceMock_FindAllMasterServices_Call) Return(_a0 model.MasterServicesOutput, _a1 error) *MasterServiceServiceMock_FindAllMasterServices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_FindAllMasterServices_Call) RunAndReturn(run func(context.Context, model.FindMasterServicesInput) (model.MasterServicesOutput, error)) *MasterServiceServiceMock_FindAllMasterServices_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllMasterServicesByUsername provides a mock function with given fields: ctx, username, input
func (_m *MasterServiceServiceMock) FindAllMasterServicesByUsername(ctx context.Context, username string, input model.FindMasterServicesInput) (model.MasterServicesOutput, error) {
	ret := _m.Called(ctx, username, input)

	if len(ret) == 0 {
		panic("no return value specified for FindAllMasterServicesByUsername")
	}

	var r0 model.MasterServicesOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FindMasterServicesInput) (model.MasterServicesOutput, error)); ok {
		return rf(ctx, username, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, model.FindMasterServicesInput) model.MasterServicesOutput); ok {
		r0 = rf(ctx, username, input)
	} else {
		r0 = ret.Get(0).(model.MasterServicesOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, model.FindMasterServicesInput) error); ok {
		r1 = rf(ctx, username, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_FindAllMasterServicesByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllMasterServicesByUsername'
type MasterServiceServiceMock_FindAllMasterServicesByUsername_Call struct {
	*mock.Call
}

// FindAllMasterServicesByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - input model.FindMasterServicesInput
func (_e *MasterServiceServiceMock_Expecter) FindAllMasterServicesByUsername(ctx interface{}, username interface{}, input interface{}) *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call {
	return &MasterServiceServiceMock_FindAllMasterServicesByUsername_Call{Call: _e.mock.On("FindAllMasterServicesByUsername", ctx, username, input)}
}

func (_c *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call) Run(run func(ctx context.Context, username string, input model.FindMasterServicesInput)) *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(model.FindMasterServicesInput))
	})
	return _c
}

func (_c *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call) Return(_a0 model.MasterServicesOutput, _a1 error) *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call) RunAndReturn(run func(context.Context, string, model.FindMasterServicesInput) (model.MasterServicesOutput, error)) *MasterServiceServiceMock_FindAllMasterServicesByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// FindAllOwnMasterServices provides a mock function with given fields: ctx, userID, input
func (_m *MasterServiceServiceMock) FindAllOwnMasterServices(ctx context.Context, userID uuid.UUID, input model.FindMasterServicesInput) (model.MasterServicesOutput, error) {
	ret := _m.Called(ctx, userID, input)

	if len(ret) == 0 {
		panic("no return value specified for FindAllOwnMasterServices")
	}

	var r0 model.MasterServicesOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.FindMasterServicesInput) (model.MasterServicesOutput, error)); ok {
		return rf(ctx, userID, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, model.FindMasterServicesInput) model.MasterServicesOutput); ok {
		r0 = rf(ctx, userID, input)
	} else {
		r0 = ret.Get(0).(model.MasterServicesOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, model.FindMasterServicesInput) error); ok {
		r1 = rf(ctx, userID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_FindAllOwnMasterServices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindAllOwnMasterServices'
type MasterServiceServiceMock_FindAllOwnMasterServices_Call struct {
	*mock.Call
}

// FindAllOwnMasterServices is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - input model.FindMasterServicesInput
func (_e *MasterServiceServiceMock_Expecter) FindAllOwnMasterServices(ctx interface{}, userID interface{}, input interface{}) *MasterServiceServiceMock_FindAllOwnMasterServices_Call {
	return &MasterServiceServiceMock_FindAllOwnMasterServices_Call{Call: _e.mock.On("FindAllOwnMasterServices", ctx, userID, input)}
}

func (_c *MasterServiceServiceMock_FindAllOwnMasterServices_Call) Run(run func(ctx context.Context, userID uuid.UUID, input model.FindMasterServicesInput)) *MasterServiceServiceMock_FindAllOwnMasterServices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(model.FindMasterServicesInput))
	})
	return _c
}

func (_c *MasterServiceServiceMock_FindAllOwnMasterServices_Call) Return(_a0 model.MasterServicesOutput, _a1 error) *MasterServiceServiceMock_FindAllOwnMasterServices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_FindAllOwnMasterServices_Call) RunAndReturn(run func(context.Context, uuid.UUID, model.FindMasterServicesInput) (model.MasterServicesOutput, error)) *MasterServiceServiceMock_FindAllOwnMasterServices_Call {
	_c.Call.Return(run)
	return _c
}

// GetMasterServiceByID provides a mock function with given fields: ctx, username, id
func (_m *MasterServiceServiceMock) GetMasterServiceByID(ctx context.Context, username string, id uuid.UUID) (model.MasterServiceOutput, error) {
	ret := _m.Called(ctx, username, id)

	if len(ret) == 0 {
		panic("no return value specified for GetMasterServiceByID")
	}

	var r0 model.MasterServiceOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) (model.MasterServiceOutput, error)); ok {
		return rf(ctx, username, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, uuid.UUID) model.MasterServiceOutput); ok {
		r0 = rf(ctx, username, id)
	} else {
		r0 = ret.Get(0).(model.MasterServiceOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, uuid.UUID) error); ok {
		r1 = rf(ctx, username, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_GetMasterServiceByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMasterServiceByID'
type MasterServiceServiceMock_GetMasterServiceByID_Call struct {
	*mock.Call
}

// GetMasterServiceByID is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - id uuid.UUID
func (_e *MasterServiceServiceMock_Expecter) GetMasterServiceByID(ctx interface{}, username interface{}, id interface{}) *MasterServiceServiceMock_GetMasterServiceByID_Call {
	return &MasterServiceServiceMock_GetMasterServiceByID_Call{Call: _e.mock.On("GetMasterServiceByID", ctx, username, id)}
}

func (_c *MasterServiceServiceMock_GetMasterServiceByID_Call) Run(run func(ctx context.Context, username string, id uuid.UUID)) *MasterServiceServiceMock_GetMasterServiceByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceServiceMock_GetMasterServiceByID_Call) Return(_a0 model.MasterServiceOutput, _a1 error) *MasterServiceServiceMock_GetMasterServiceByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_GetMasterServiceByID_Call) RunAndReturn(run func(context.Context, string, uuid.UUID) (model.MasterServiceOutput, error)) *MasterServiceServiceMock_GetMasterServiceByID_Call {
	_c.Call.Return(run)
	return _c
}

// GetOwnMasterServiceByID provides a mock function with given fields: ctx, userID, id
func (_m *MasterServiceServiceMock) GetOwnMasterServiceByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (model.MasterServiceOutput, error) {
	ret := _m.Called(ctx, userID, id)

	if len(ret) == 0 {
		panic("no return value specified for GetOwnMasterServiceByID")
	}

	var r0 model.MasterServiceOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (model.MasterServiceOutput, error)); ok {
		return rf(ctx, userID, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) model.MasterServiceOutput); ok {
		r0 = rf(ctx, userID, id)
	} else {
		r0 = ret.Get(0).(model.MasterServiceOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, userID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_GetOwnMasterServiceByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOwnMasterServiceByID'
type MasterServiceServiceMock_GetOwnMasterServiceByID_Call struct {
	*mock.Call
}

// GetOwnMasterServiceByID is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - id uuid.UUID
func (_e *MasterServiceServiceMock_Expecter) GetOwnMasterServiceByID(ctx interface{}, userID interface{}, id interface{}) *MasterServiceServiceMock_GetOwnMasterServiceByID_Call {
	return &MasterServiceServiceMock_GetOwnMasterServiceByID_Call{Call: _e.mock.On("GetOwnMasterServiceByID", ctx, userID, id)}
}

func (_c *MasterServiceServiceMock_GetOwnMasterServiceByID_Call) Run(run func(ctx context.Context, userID uuid.UUID, id uuid.UUID)) *MasterServiceServiceMock_GetOwnMasterServiceByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceServiceMock_GetOwnMasterServiceByID_Call) Return(_a0 model.MasterServiceOutput, _a1 error) *MasterServiceServiceMock_GetOwnMasterServiceByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_GetOwnMasterServiceByID_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (model.MasterServiceOutput, error)) *MasterServiceServiceMock_GetOwnMasterServiceByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMasterService provides a mock function with given fields: ctx, userID, masterServiceID, input
func (_m *MasterServiceServiceMock) UpdateMasterService(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID, input model.UpdateMasterServiceInput) (model.MasterServiceOutput, error) {
	ret := _m.Called(ctx, userID, masterServiceID, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMasterService")
	}

	var r0 model.MasterServiceOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, model.UpdateMasterServiceInput) (model.MasterServiceOutput, error)); ok {
		return rf(ctx, userID, masterServiceID, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID, model.UpdateMasterServiceInput) model.MasterServiceOutput); ok {
		r0 = rf(ctx, userID, masterServiceID, input)
	} else {
		r0 = ret.Get(0).(model.MasterServiceOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID, model.UpdateMasterServiceInput) error); ok {
		r1 = rf(ctx, userID, masterServiceID, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceServiceMock_UpdateMasterService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMasterService'
type MasterServiceServiceMock_UpdateMasterService_Call struct {
	*mock.Call
}

// UpdateMasterService is a helper method to define mock.On call
//   - ctx context.Context
//   - userID uuid.UUID
//   - masterServiceID uuid.UUID
//   - input model.UpdateMasterServiceInput
func (_e *MasterServiceServiceMock_Expecter) UpdateMasterService(ctx interface{}, userID interface{}, masterServiceID interface{}, input interface{}) *MasterServiceServiceMock_UpdateMasterService_Call {
	return &MasterServiceServiceMock_UpdateMasterService_Call{Call: _e.mock.On("UpdateMasterService", ctx, userID, masterServiceID, input)}
}

func (_c *MasterServiceServiceMock_UpdateMasterService_Call) Run(run func(ctx context.Context, userID uuid.UUID, masterServiceID uuid.UUID, input model.UpdateMasterServiceInput)) *MasterServiceServiceMock_UpdateMasterService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID), args[3].(model.UpdateMasterServiceInput))
	})
	return _c
}

func (_c *MasterServiceServiceMock_UpdateMasterService_Call) Return(_a0 model.MasterServiceOutput, _a1 error) *MasterServiceServiceMock_UpdateMasterService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceServiceMock_UpdateMasterService_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID, model.UpdateMasterServiceInput) (model.MasterServiceOutput, error)) *MasterServiceServiceMock_UpdateMasterService_Call {
	_c.Call.Return(run)
	return _c
}

// NewMasterServiceServiceMock creates a new instance of MasterServiceServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMasterServiceServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MasterServiceServiceMock {
	mock := &MasterServiceServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
