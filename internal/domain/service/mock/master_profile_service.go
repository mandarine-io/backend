// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	dto "github.com/mandarine-io/Backend/internal/domain/dto"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// MasterProfileServiceMock is an autogenerated mock type for the MasterProfileService type
type MasterProfileServiceMock struct {
	mock.Mock
}

type MasterProfileServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MasterProfileServiceMock) EXPECT() *MasterProfileServiceMock_Expecter {
	return &MasterProfileServiceMock_Expecter{mock: &_m.Mock}
}

// CreateMasterProfile provides a mock function with given fields: ctx, id, input
func (_m *MasterProfileServiceMock) CreateMasterProfile(ctx context.Context, id uuid.UUID, input dto.CreateMasterProfileInput) (dto.OwnMasterProfileOutput, error) {
	ret := _m.Called(ctx, id, input)

	if len(ret) == 0 {
		panic("no return value specified for CreateMasterProfile")
	}

	var r0 dto.OwnMasterProfileOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.CreateMasterProfileInput) (dto.OwnMasterProfileOutput, error)); ok {
		return rf(ctx, id, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.CreateMasterProfileInput) dto.OwnMasterProfileOutput); ok {
		r0 = rf(ctx, id, input)
	} else {
		r0 = ret.Get(0).(dto.OwnMasterProfileOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, dto.CreateMasterProfileInput) error); ok {
		r1 = rf(ctx, id, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterProfileServiceMock_CreateMasterProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMasterProfile'
type MasterProfileServiceMock_CreateMasterProfile_Call struct {
	*mock.Call
}

// CreateMasterProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - input dto.CreateMasterProfileInput
func (_e *MasterProfileServiceMock_Expecter) CreateMasterProfile(ctx interface{}, id interface{}, input interface{}) *MasterProfileServiceMock_CreateMasterProfile_Call {
	return &MasterProfileServiceMock_CreateMasterProfile_Call{Call: _e.mock.On("CreateMasterProfile", ctx, id, input)}
}

func (_c *MasterProfileServiceMock_CreateMasterProfile_Call) Run(run func(ctx context.Context, id uuid.UUID, input dto.CreateMasterProfileInput)) *MasterProfileServiceMock_CreateMasterProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(dto.CreateMasterProfileInput))
	})
	return _c
}

func (_c *MasterProfileServiceMock_CreateMasterProfile_Call) Return(_a0 dto.OwnMasterProfileOutput, _a1 error) *MasterProfileServiceMock_CreateMasterProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterProfileServiceMock_CreateMasterProfile_Call) RunAndReturn(run func(context.Context, uuid.UUID, dto.CreateMasterProfileInput) (dto.OwnMasterProfileOutput, error)) *MasterProfileServiceMock_CreateMasterProfile_Call {
	_c.Call.Return(run)
	return _c
}

// FindMasterProfiles provides a mock function with given fields: ctx, input
func (_m *MasterProfileServiceMock) FindMasterProfiles(ctx context.Context, input dto.FindMasterProfilesInput) (dto.MasterProfilesOutput, error) {
	ret := _m.Called(ctx, input)

	if len(ret) == 0 {
		panic("no return value specified for FindMasterProfiles")
	}

	var r0 dto.MasterProfilesOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, dto.FindMasterProfilesInput) (dto.MasterProfilesOutput, error)); ok {
		return rf(ctx, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, dto.FindMasterProfilesInput) dto.MasterProfilesOutput); ok {
		r0 = rf(ctx, input)
	} else {
		r0 = ret.Get(0).(dto.MasterProfilesOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, dto.FindMasterProfilesInput) error); ok {
		r1 = rf(ctx, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterProfileServiceMock_FindMasterProfiles_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMasterProfiles'
type MasterProfileServiceMock_FindMasterProfiles_Call struct {
	*mock.Call
}

// FindMasterProfiles is a helper method to define mock.On call
//   - ctx context.Context
//   - input dto.FindMasterProfilesInput
func (_e *MasterProfileServiceMock_Expecter) FindMasterProfiles(ctx interface{}, input interface{}) *MasterProfileServiceMock_FindMasterProfiles_Call {
	return &MasterProfileServiceMock_FindMasterProfiles_Call{Call: _e.mock.On("FindMasterProfiles", ctx, input)}
}

func (_c *MasterProfileServiceMock_FindMasterProfiles_Call) Run(run func(ctx context.Context, input dto.FindMasterProfilesInput)) *MasterProfileServiceMock_FindMasterProfiles_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.FindMasterProfilesInput))
	})
	return _c
}

func (_c *MasterProfileServiceMock_FindMasterProfiles_Call) Return(_a0 dto.MasterProfilesOutput, _a1 error) *MasterProfileServiceMock_FindMasterProfiles_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterProfileServiceMock_FindMasterProfiles_Call) RunAndReturn(run func(context.Context, dto.FindMasterProfilesInput) (dto.MasterProfilesOutput, error)) *MasterProfileServiceMock_FindMasterProfiles_Call {
	_c.Call.Return(run)
	return _c
}

// GetMasterProfileByUsername provides a mock function with given fields: ctx, username
func (_m *MasterProfileServiceMock) GetMasterProfileByUsername(ctx context.Context, username string) (dto.MasterProfileOutput, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for GetMasterProfileByUsername")
	}

	var r0 dto.MasterProfileOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (dto.MasterProfileOutput, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) dto.MasterProfileOutput); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(dto.MasterProfileOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterProfileServiceMock_GetMasterProfileByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMasterProfileByUsername'
type MasterProfileServiceMock_GetMasterProfileByUsername_Call struct {
	*mock.Call
}

// GetMasterProfileByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *MasterProfileServiceMock_Expecter) GetMasterProfileByUsername(ctx interface{}, username interface{}) *MasterProfileServiceMock_GetMasterProfileByUsername_Call {
	return &MasterProfileServiceMock_GetMasterProfileByUsername_Call{Call: _e.mock.On("GetMasterProfileByUsername", ctx, username)}
}

func (_c *MasterProfileServiceMock_GetMasterProfileByUsername_Call) Run(run func(ctx context.Context, username string)) *MasterProfileServiceMock_GetMasterProfileByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MasterProfileServiceMock_GetMasterProfileByUsername_Call) Return(_a0 dto.MasterProfileOutput, _a1 error) *MasterProfileServiceMock_GetMasterProfileByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterProfileServiceMock_GetMasterProfileByUsername_Call) RunAndReturn(run func(context.Context, string) (dto.MasterProfileOutput, error)) *MasterProfileServiceMock_GetMasterProfileByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// GetOwnMasterProfile provides a mock function with given fields: ctx, id
func (_m *MasterProfileServiceMock) GetOwnMasterProfile(ctx context.Context, id uuid.UUID) (dto.OwnMasterProfileOutput, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetOwnMasterProfile")
	}

	var r0 dto.OwnMasterProfileOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (dto.OwnMasterProfileOutput, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) dto.OwnMasterProfileOutput); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(dto.OwnMasterProfileOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterProfileServiceMock_GetOwnMasterProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOwnMasterProfile'
type MasterProfileServiceMock_GetOwnMasterProfile_Call struct {
	*mock.Call
}

// GetOwnMasterProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *MasterProfileServiceMock_Expecter) GetOwnMasterProfile(ctx interface{}, id interface{}) *MasterProfileServiceMock_GetOwnMasterProfile_Call {
	return &MasterProfileServiceMock_GetOwnMasterProfile_Call{Call: _e.mock.On("GetOwnMasterProfile", ctx, id)}
}

func (_c *MasterProfileServiceMock_GetOwnMasterProfile_Call) Run(run func(ctx context.Context, id uuid.UUID)) *MasterProfileServiceMock_GetOwnMasterProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MasterProfileServiceMock_GetOwnMasterProfile_Call) Return(_a0 dto.OwnMasterProfileOutput, _a1 error) *MasterProfileServiceMock_GetOwnMasterProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterProfileServiceMock_GetOwnMasterProfile_Call) RunAndReturn(run func(context.Context, uuid.UUID) (dto.OwnMasterProfileOutput, error)) *MasterProfileServiceMock_GetOwnMasterProfile_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMasterProfile provides a mock function with given fields: ctx, id, input
func (_m *MasterProfileServiceMock) UpdateMasterProfile(ctx context.Context, id uuid.UUID, input dto.UpdateMasterProfileInput) (dto.OwnMasterProfileOutput, error) {
	ret := _m.Called(ctx, id, input)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMasterProfile")
	}

	var r0 dto.OwnMasterProfileOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.UpdateMasterProfileInput) (dto.OwnMasterProfileOutput, error)); ok {
		return rf(ctx, id, input)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, dto.UpdateMasterProfileInput) dto.OwnMasterProfileOutput); ok {
		r0 = rf(ctx, id, input)
	} else {
		r0 = ret.Get(0).(dto.OwnMasterProfileOutput)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, dto.UpdateMasterProfileInput) error); ok {
		r1 = rf(ctx, id, input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterProfileServiceMock_UpdateMasterProfile_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMasterProfile'
type MasterProfileServiceMock_UpdateMasterProfile_Call struct {
	*mock.Call
}

// UpdateMasterProfile is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - input dto.UpdateMasterProfileInput
func (_e *MasterProfileServiceMock_Expecter) UpdateMasterProfile(ctx interface{}, id interface{}, input interface{}) *MasterProfileServiceMock_UpdateMasterProfile_Call {
	return &MasterProfileServiceMock_UpdateMasterProfile_Call{Call: _e.mock.On("UpdateMasterProfile", ctx, id, input)}
}

func (_c *MasterProfileServiceMock_UpdateMasterProfile_Call) Run(run func(ctx context.Context, id uuid.UUID, input dto.UpdateMasterProfileInput)) *MasterProfileServiceMock_UpdateMasterProfile_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(dto.UpdateMasterProfileInput))
	})
	return _c
}

func (_c *MasterProfileServiceMock_UpdateMasterProfile_Call) Return(_a0 dto.OwnMasterProfileOutput, _a1 error) *MasterProfileServiceMock_UpdateMasterProfile_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterProfileServiceMock_UpdateMasterProfile_Call) RunAndReturn(run func(context.Context, uuid.UUID, dto.UpdateMasterProfileInput) (dto.OwnMasterProfileOutput, error)) *MasterProfileServiceMock_UpdateMasterProfile_Call {
	_c.Call.Return(run)
	return _c
}

// NewMasterProfileServiceMock creates a new instance of MasterProfileServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMasterProfileServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MasterProfileServiceMock {
	mock := &MasterProfileServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
