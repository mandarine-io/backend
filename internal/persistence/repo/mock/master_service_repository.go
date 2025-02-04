// Code generated by mockery v2.51.1. DO NOT EDIT.

package mock

import (
	context "context"

	entity "github.com/mandarine-io/backend/internal/persistence/entity"
	decimal "github.com/shopspring/decimal"

	mock "github.com/stretchr/testify/mock"

	repo "github.com/mandarine-io/backend/internal/persistence/repo"

	time "time"

	uuid "github.com/google/uuid"
)

// MasterServiceRepositoryMock is an autogenerated mock type for the MasterServiceRepository type
type MasterServiceRepositoryMock struct {
	mock.Mock
}

type MasterServiceRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *MasterServiceRepositoryMock) EXPECT() *MasterServiceRepositoryMock_Expecter {
	return &MasterServiceRepositoryMock_Expecter{mock: &_m.Mock}
}

// CountMasterServices provides a mock function with given fields: ctx, scopes
func (_m *MasterServiceRepositoryMock) CountMasterServices(ctx context.Context, scopes ...repo.Scope) (int64, error) {
	_va := make([]interface{}, len(scopes))
	for _i := range scopes {
		_va[_i] = scopes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CountMasterServices")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.Scope) (int64, error)); ok {
		return rf(ctx, scopes...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.Scope) int64); ok {
		r0 = rf(ctx, scopes...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...repo.Scope) error); ok {
		r1 = rf(ctx, scopes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_CountMasterServices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountMasterServices'
type MasterServiceRepositoryMock_CountMasterServices_Call struct {
	*mock.Call
}

// CountMasterServices is a helper method to define mock.On call
//   - ctx context.Context
//   - scopes ...repo.Scope
func (_e *MasterServiceRepositoryMock_Expecter) CountMasterServices(ctx interface{}, scopes ...interface{}) *MasterServiceRepositoryMock_CountMasterServices_Call {
	return &MasterServiceRepositoryMock_CountMasterServices_Call{Call: _e.mock.On("CountMasterServices",
		append([]interface{}{ctx}, scopes...)...)}
}

func (_c *MasterServiceRepositoryMock_CountMasterServices_Call) Run(run func(ctx context.Context, scopes ...repo.Scope)) *MasterServiceRepositoryMock_CountMasterServices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]repo.Scope, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(repo.Scope)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_CountMasterServices_Call) Return(_a0 int64, _a1 error) *MasterServiceRepositoryMock_CountMasterServices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_CountMasterServices_Call) RunAndReturn(run func(context.Context, ...repo.Scope) (int64, error)) *MasterServiceRepositoryMock_CountMasterServices_Call {
	_c.Call.Return(run)
	return _c
}

// CountMasterServicesByMasterProfileID provides a mock function with given fields: ctx, masterProfileID, scopes
func (_m *MasterServiceRepositoryMock) CountMasterServicesByMasterProfileID(ctx context.Context, masterProfileID uuid.UUID, scopes ...repo.Scope) (int64, error) {
	_va := make([]interface{}, len(scopes))
	for _i := range scopes {
		_va[_i] = scopes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, masterProfileID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for CountMasterServicesByMasterProfileID")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, ...repo.Scope) (int64, error)); ok {
		return rf(ctx, masterProfileID, scopes...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, ...repo.Scope) int64); ok {
		r0 = rf(ctx, masterProfileID, scopes...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, ...repo.Scope) error); ok {
		r1 = rf(ctx, masterProfileID, scopes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountMasterServicesByMasterProfileID'
type MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call struct {
	*mock.Call
}

// CountMasterServicesByMasterProfileID is a helper method to define mock.On call
//   - ctx context.Context
//   - masterProfileID uuid.UUID
//   - scopes ...repo.Scope
func (_e *MasterServiceRepositoryMock_Expecter) CountMasterServicesByMasterProfileID(ctx interface{}, masterProfileID interface{}, scopes ...interface{}) *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call {
	return &MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call{Call: _e.mock.On("CountMasterServicesByMasterProfileID",
		append([]interface{}{ctx, masterProfileID}, scopes...)...)}
}

func (_c *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call) Run(run func(ctx context.Context, masterProfileID uuid.UUID, scopes ...repo.Scope)) *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]repo.Scope, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(repo.Scope)
			}
		}
		run(args[0].(context.Context), args[1].(uuid.UUID), variadicArgs...)
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call) Return(_a0 int64, _a1 error) *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call) RunAndReturn(run func(context.Context, uuid.UUID, ...repo.Scope) (int64, error)) *MasterServiceRepositoryMock_CountMasterServicesByMasterProfileID_Call {
	_c.Call.Return(run)
	return _c
}

// CreateMasterService provides a mock function with given fields: ctx, masterService
func (_m *MasterServiceRepositoryMock) CreateMasterService(ctx context.Context, masterService *entity.MasterService) (*entity.MasterService, error) {
	ret := _m.Called(ctx, masterService)

	if len(ret) == 0 {
		panic("no return value specified for CreateMasterService")
	}

	var r0 *entity.MasterService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.MasterService) (*entity.MasterService, error)); ok {
		return rf(ctx, masterService)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.MasterService) *entity.MasterService); ok {
		r0 = rf(ctx, masterService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.MasterService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.MasterService) error); ok {
		r1 = rf(ctx, masterService)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_CreateMasterService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateMasterService'
type MasterServiceRepositoryMock_CreateMasterService_Call struct {
	*mock.Call
}

// CreateMasterService is a helper method to define mock.On call
//   - ctx context.Context
//   - masterService *entity.MasterService
func (_e *MasterServiceRepositoryMock_Expecter) CreateMasterService(ctx interface{}, masterService interface{}) *MasterServiceRepositoryMock_CreateMasterService_Call {
	return &MasterServiceRepositoryMock_CreateMasterService_Call{Call: _e.mock.On("CreateMasterService", ctx, masterService)}
}

func (_c *MasterServiceRepositoryMock_CreateMasterService_Call) Run(run func(ctx context.Context, masterService *entity.MasterService)) *MasterServiceRepositoryMock_CreateMasterService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.MasterService))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_CreateMasterService_Call) Return(_a0 *entity.MasterService, _a1 error) *MasterServiceRepositoryMock_CreateMasterService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_CreateMasterService_Call) RunAndReturn(run func(context.Context, *entity.MasterService) (*entity.MasterService, error)) *MasterServiceRepositoryMock_CreateMasterService_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMasterServiceByID provides a mock function with given fields: ctx, masterProfileID, id
func (_m *MasterServiceRepositoryMock) DeleteMasterServiceByID(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID) error {
	ret := _m.Called(ctx, masterProfileID, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMasterServiceByID")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r0 = rf(ctx, masterProfileID, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MasterServiceRepositoryMock_DeleteMasterServiceByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMasterServiceByID'
type MasterServiceRepositoryMock_DeleteMasterServiceByID_Call struct {
	*mock.Call
}

// DeleteMasterServiceByID is a helper method to define mock.On call
//   - ctx context.Context
//   - masterProfileID uuid.UUID
//   - id uuid.UUID
func (_e *MasterServiceRepositoryMock_Expecter) DeleteMasterServiceByID(ctx interface{}, masterProfileID interface{}, id interface{}) *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call {
	return &MasterServiceRepositoryMock_DeleteMasterServiceByID_Call{Call: _e.mock.On("DeleteMasterServiceByID", ctx, masterProfileID, id)}
}

func (_c *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call) Run(run func(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID)) *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call) Return(_a0 error) *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) error) *MasterServiceRepositoryMock_DeleteMasterServiceByID_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsMasterServiceByMasterID provides a mock function with given fields: ctx, masterID
func (_m *MasterServiceRepositoryMock) ExistsMasterServiceByMasterID(ctx context.Context, masterID uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, masterID)

	if len(ret) == 0 {
		panic("no return value specified for ExistsMasterServiceByMasterID")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (bool, error)); ok {
		return rf(ctx, masterID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) bool); ok {
		r0 = rf(ctx, masterID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, masterID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsMasterServiceByMasterID'
type MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call struct {
	*mock.Call
}

// ExistsMasterServiceByMasterID is a helper method to define mock.On call
//   - ctx context.Context
//   - masterID uuid.UUID
func (_e *MasterServiceRepositoryMock_Expecter) ExistsMasterServiceByMasterID(ctx interface{}, masterID interface{}) *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call {
	return &MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call{Call: _e.mock.On("ExistsMasterServiceByMasterID", ctx, masterID)}
}

func (_c *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call) Run(run func(ctx context.Context, masterID uuid.UUID)) *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call) Return(_a0 bool, _a1 error) *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call) RunAndReturn(run func(context.Context, uuid.UUID) (bool, error)) *MasterServiceRepositoryMock_ExistsMasterServiceByMasterID_Call {
	_c.Call.Return(run)
	return _c
}

// FindMasterServiceByID provides a mock function with given fields: ctx, masterProfileID, id
func (_m *MasterServiceRepositoryMock) FindMasterServiceByID(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID) (*entity.MasterService, error) {
	ret := _m.Called(ctx, masterProfileID, id)

	if len(ret) == 0 {
		panic("no return value specified for FindMasterServiceByID")
	}

	var r0 *entity.MasterService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) (*entity.MasterService, error)); ok {
		return rf(ctx, masterProfileID, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, uuid.UUID) *entity.MasterService); ok {
		r0 = rf(ctx, masterProfileID, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.MasterService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, uuid.UUID) error); ok {
		r1 = rf(ctx, masterProfileID, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_FindMasterServiceByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMasterServiceByID'
type MasterServiceRepositoryMock_FindMasterServiceByID_Call struct {
	*mock.Call
}

// FindMasterServiceByID is a helper method to define mock.On call
//   - ctx context.Context
//   - masterProfileID uuid.UUID
//   - id uuid.UUID
func (_e *MasterServiceRepositoryMock_Expecter) FindMasterServiceByID(ctx interface{}, masterProfileID interface{}, id interface{}) *MasterServiceRepositoryMock_FindMasterServiceByID_Call {
	return &MasterServiceRepositoryMock_FindMasterServiceByID_Call{Call: _e.mock.On("FindMasterServiceByID", ctx, masterProfileID, id)}
}

func (_c *MasterServiceRepositoryMock_FindMasterServiceByID_Call) Run(run func(ctx context.Context, masterProfileID uuid.UUID, id uuid.UUID)) *MasterServiceRepositoryMock_FindMasterServiceByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(uuid.UUID))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServiceByID_Call) Return(_a0 *entity.MasterService, _a1 error) *MasterServiceRepositoryMock_FindMasterServiceByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServiceByID_Call) RunAndReturn(run func(context.Context, uuid.UUID, uuid.UUID) (*entity.MasterService, error)) *MasterServiceRepositoryMock_FindMasterServiceByID_Call {
	_c.Call.Return(run)
	return _c
}

// FindMasterServices provides a mock function with given fields: ctx, scopes
func (_m *MasterServiceRepositoryMock) FindMasterServices(ctx context.Context, scopes ...repo.Scope) ([]*entity.MasterService, error) {
	_va := make([]interface{}, len(scopes))
	for _i := range scopes {
		_va[_i] = scopes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindMasterServices")
	}

	var r0 []*entity.MasterService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.Scope) ([]*entity.MasterService, error)); ok {
		return rf(ctx, scopes...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...repo.Scope) []*entity.MasterService); ok {
		r0 = rf(ctx, scopes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.MasterService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...repo.Scope) error); ok {
		r1 = rf(ctx, scopes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_FindMasterServices_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMasterServices'
type MasterServiceRepositoryMock_FindMasterServices_Call struct {
	*mock.Call
}

// FindMasterServices is a helper method to define mock.On call
//   - ctx context.Context
//   - scopes ...repo.Scope
func (_e *MasterServiceRepositoryMock_Expecter) FindMasterServices(ctx interface{}, scopes ...interface{}) *MasterServiceRepositoryMock_FindMasterServices_Call {
	return &MasterServiceRepositoryMock_FindMasterServices_Call{Call: _e.mock.On("FindMasterServices",
		append([]interface{}{ctx}, scopes...)...)}
}

func (_c *MasterServiceRepositoryMock_FindMasterServices_Call) Run(run func(ctx context.Context, scopes ...repo.Scope)) *MasterServiceRepositoryMock_FindMasterServices_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]repo.Scope, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(repo.Scope)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServices_Call) Return(_a0 []*entity.MasterService, _a1 error) *MasterServiceRepositoryMock_FindMasterServices_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServices_Call) RunAndReturn(run func(context.Context, ...repo.Scope) ([]*entity.MasterService, error)) *MasterServiceRepositoryMock_FindMasterServices_Call {
	_c.Call.Return(run)
	return _c
}

// FindMasterServicesByMasterProfileID provides a mock function with given fields: ctx, masterProfileID, scopes
func (_m *MasterServiceRepositoryMock) FindMasterServicesByMasterProfileID(ctx context.Context, masterProfileID uuid.UUID, scopes ...repo.Scope) ([]*entity.MasterService, error) {
	_va := make([]interface{}, len(scopes))
	for _i := range scopes {
		_va[_i] = scopes[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, masterProfileID)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindMasterServicesByMasterProfileID")
	}

	var r0 []*entity.MasterService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, ...repo.Scope) ([]*entity.MasterService, error)); ok {
		return rf(ctx, masterProfileID, scopes...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, ...repo.Scope) []*entity.MasterService); ok {
		r0 = rf(ctx, masterProfileID, scopes...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.MasterService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, ...repo.Scope) error); ok {
		r1 = rf(ctx, masterProfileID, scopes...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindMasterServicesByMasterProfileID'
type MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call struct {
	*mock.Call
}

// FindMasterServicesByMasterProfileID is a helper method to define mock.On call
//   - ctx context.Context
//   - masterProfileID uuid.UUID
//   - scopes ...repo.Scope
func (_e *MasterServiceRepositoryMock_Expecter) FindMasterServicesByMasterProfileID(ctx interface{}, masterProfileID interface{}, scopes ...interface{}) *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call {
	return &MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call{Call: _e.mock.On("FindMasterServicesByMasterProfileID",
		append([]interface{}{ctx, masterProfileID}, scopes...)...)}
}

func (_c *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call) Run(run func(ctx context.Context, masterProfileID uuid.UUID, scopes ...repo.Scope)) *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]repo.Scope, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(repo.Scope)
			}
		}
		run(args[0].(context.Context), args[1].(uuid.UUID), variadicArgs...)
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call) Return(_a0 []*entity.MasterService, _a1 error) *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call) RunAndReturn(run func(context.Context, uuid.UUID, ...repo.Scope) ([]*entity.MasterService, error)) *MasterServiceRepositoryMock_FindMasterServicesByMasterProfileID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMasterService provides a mock function with given fields: ctx, masterService
func (_m *MasterServiceRepositoryMock) UpdateMasterService(ctx context.Context, masterService *entity.MasterService) (*entity.MasterService, error) {
	ret := _m.Called(ctx, masterService)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMasterService")
	}

	var r0 *entity.MasterService
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *entity.MasterService) (*entity.MasterService, error)); ok {
		return rf(ctx, masterService)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *entity.MasterService) *entity.MasterService); ok {
		r0 = rf(ctx, masterService)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.MasterService)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *entity.MasterService) error); ok {
		r1 = rf(ctx, masterService)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MasterServiceRepositoryMock_UpdateMasterService_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMasterService'
type MasterServiceRepositoryMock_UpdateMasterService_Call struct {
	*mock.Call
}

// UpdateMasterService is a helper method to define mock.On call
//   - ctx context.Context
//   - masterService *entity.MasterService
func (_e *MasterServiceRepositoryMock_Expecter) UpdateMasterService(ctx interface{}, masterService interface{}) *MasterServiceRepositoryMock_UpdateMasterService_Call {
	return &MasterServiceRepositoryMock_UpdateMasterService_Call{Call: _e.mock.On("UpdateMasterService", ctx, masterService)}
}

func (_c *MasterServiceRepositoryMock_UpdateMasterService_Call) Run(run func(ctx context.Context, masterService *entity.MasterService)) *MasterServiceRepositoryMock_UpdateMasterService_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entity.MasterService))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_UpdateMasterService_Call) Return(_a0 *entity.MasterService, _a1 error) *MasterServiceRepositoryMock_UpdateMasterService_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MasterServiceRepositoryMock_UpdateMasterService_Call) RunAndReturn(run func(context.Context, *entity.MasterService) (*entity.MasterService, error)) *MasterServiceRepositoryMock_UpdateMasterService_Call {
	_c.Call.Return(run)
	return _c
}

// WithMaxIntervalFilter provides a mock function with given fields: maxDuration
func (_m *MasterServiceRepositoryMock) WithMaxIntervalFilter(maxDuration time.Duration) repo.Scope {
	ret := _m.Called(maxDuration)

	if len(ret) == 0 {
		panic("no return value specified for WithMaxIntervalFilter")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(time.Duration) repo.Scope); ok {
		r0 = rf(maxDuration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithMaxIntervalFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMaxIntervalFilter'
type MasterServiceRepositoryMock_WithMaxIntervalFilter_Call struct {
	*mock.Call
}

// WithMaxIntervalFilter is a helper method to define mock.On call
//   - maxDuration time.Duration
func (_e *MasterServiceRepositoryMock_Expecter) WithMaxIntervalFilter(maxDuration interface{}) *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call {
	return &MasterServiceRepositoryMock_WithMaxIntervalFilter_Call{Call: _e.mock.On("WithMaxIntervalFilter", maxDuration)}
}

func (_c *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call) Run(run func(maxDuration time.Duration)) *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call) RunAndReturn(run func(time.Duration) repo.Scope) *MasterServiceRepositoryMock_WithMaxIntervalFilter_Call {
	_c.Call.Return(run)
	return _c
}

// WithMaxPriceFilter provides a mock function with given fields: maxPrice
func (_m *MasterServiceRepositoryMock) WithMaxPriceFilter(maxPrice decimal.Decimal) repo.Scope {
	ret := _m.Called(maxPrice)

	if len(ret) == 0 {
		panic("no return value specified for WithMaxPriceFilter")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(decimal.Decimal) repo.Scope); ok {
		r0 = rf(maxPrice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithMaxPriceFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMaxPriceFilter'
type MasterServiceRepositoryMock_WithMaxPriceFilter_Call struct {
	*mock.Call
}

// WithMaxPriceFilter is a helper method to define mock.On call
//   - maxPrice decimal.Decimal
func (_e *MasterServiceRepositoryMock_Expecter) WithMaxPriceFilter(maxPrice interface{}) *MasterServiceRepositoryMock_WithMaxPriceFilter_Call {
	return &MasterServiceRepositoryMock_WithMaxPriceFilter_Call{Call: _e.mock.On("WithMaxPriceFilter", maxPrice)}
}

func (_c *MasterServiceRepositoryMock_WithMaxPriceFilter_Call) Run(run func(maxPrice decimal.Decimal)) *MasterServiceRepositoryMock_WithMaxPriceFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(decimal.Decimal))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMaxPriceFilter_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithMaxPriceFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMaxPriceFilter_Call) RunAndReturn(run func(decimal.Decimal) repo.Scope) *MasterServiceRepositoryMock_WithMaxPriceFilter_Call {
	_c.Call.Return(run)
	return _c
}

// WithMinIntervalFilter provides a mock function with given fields: minDuration
func (_m *MasterServiceRepositoryMock) WithMinIntervalFilter(minDuration time.Duration) repo.Scope {
	ret := _m.Called(minDuration)

	if len(ret) == 0 {
		panic("no return value specified for WithMinIntervalFilter")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(time.Duration) repo.Scope); ok {
		r0 = rf(minDuration)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithMinIntervalFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMinIntervalFilter'
type MasterServiceRepositoryMock_WithMinIntervalFilter_Call struct {
	*mock.Call
}

// WithMinIntervalFilter is a helper method to define mock.On call
//   - minDuration time.Duration
func (_e *MasterServiceRepositoryMock_Expecter) WithMinIntervalFilter(minDuration interface{}) *MasterServiceRepositoryMock_WithMinIntervalFilter_Call {
	return &MasterServiceRepositoryMock_WithMinIntervalFilter_Call{Call: _e.mock.On("WithMinIntervalFilter", minDuration)}
}

func (_c *MasterServiceRepositoryMock_WithMinIntervalFilter_Call) Run(run func(minDuration time.Duration)) *MasterServiceRepositoryMock_WithMinIntervalFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(time.Duration))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMinIntervalFilter_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithMinIntervalFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMinIntervalFilter_Call) RunAndReturn(run func(time.Duration) repo.Scope) *MasterServiceRepositoryMock_WithMinIntervalFilter_Call {
	_c.Call.Return(run)
	return _c
}

// WithMinPriceFilter provides a mock function with given fields: minPrice
func (_m *MasterServiceRepositoryMock) WithMinPriceFilter(minPrice decimal.Decimal) repo.Scope {
	ret := _m.Called(minPrice)

	if len(ret) == 0 {
		panic("no return value specified for WithMinPriceFilter")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(decimal.Decimal) repo.Scope); ok {
		r0 = rf(minPrice)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithMinPriceFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithMinPriceFilter'
type MasterServiceRepositoryMock_WithMinPriceFilter_Call struct {
	*mock.Call
}

// WithMinPriceFilter is a helper method to define mock.On call
//   - minPrice decimal.Decimal
func (_e *MasterServiceRepositoryMock_Expecter) WithMinPriceFilter(minPrice interface{}) *MasterServiceRepositoryMock_WithMinPriceFilter_Call {
	return &MasterServiceRepositoryMock_WithMinPriceFilter_Call{Call: _e.mock.On("WithMinPriceFilter", minPrice)}
}

func (_c *MasterServiceRepositoryMock_WithMinPriceFilter_Call) Run(run func(minPrice decimal.Decimal)) *MasterServiceRepositoryMock_WithMinPriceFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(decimal.Decimal))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMinPriceFilter_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithMinPriceFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithMinPriceFilter_Call) RunAndReturn(run func(decimal.Decimal) repo.Scope) *MasterServiceRepositoryMock_WithMinPriceFilter_Call {
	_c.Call.Return(run)
	return _c
}

// WithNameFilter provides a mock function with given fields: name
func (_m *MasterServiceRepositoryMock) WithNameFilter(name string) repo.Scope {
	ret := _m.Called(name)

	if len(ret) == 0 {
		panic("no return value specified for WithNameFilter")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(string) repo.Scope); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithNameFilter_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithNameFilter'
type MasterServiceRepositoryMock_WithNameFilter_Call struct {
	*mock.Call
}

// WithNameFilter is a helper method to define mock.On call
//   - name string
func (_e *MasterServiceRepositoryMock_Expecter) WithNameFilter(name interface{}) *MasterServiceRepositoryMock_WithNameFilter_Call {
	return &MasterServiceRepositoryMock_WithNameFilter_Call{Call: _e.mock.On("WithNameFilter", name)}
}

func (_c *MasterServiceRepositoryMock_WithNameFilter_Call) Run(run func(name string)) *MasterServiceRepositoryMock_WithNameFilter_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithNameFilter_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithNameFilter_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithNameFilter_Call) RunAndReturn(run func(string) repo.Scope) *MasterServiceRepositoryMock_WithNameFilter_Call {
	_c.Call.Return(run)
	return _c
}

// WithPagination provides a mock function with given fields: page, pageSize
func (_m *MasterServiceRepositoryMock) WithPagination(page int, pageSize int) repo.Scope {
	ret := _m.Called(page, pageSize)

	if len(ret) == 0 {
		panic("no return value specified for WithPagination")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(int, int) repo.Scope); ok {
		r0 = rf(page, pageSize)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithPagination_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithPagination'
type MasterServiceRepositoryMock_WithPagination_Call struct {
	*mock.Call
}

// WithPagination is a helper method to define mock.On call
//   - page int
//   - pageSize int
func (_e *MasterServiceRepositoryMock_Expecter) WithPagination(page interface{}, pageSize interface{}) *MasterServiceRepositoryMock_WithPagination_Call {
	return &MasterServiceRepositoryMock_WithPagination_Call{Call: _e.mock.On("WithPagination", page, pageSize)}
}

func (_c *MasterServiceRepositoryMock_WithPagination_Call) Run(run func(page int, pageSize int)) *MasterServiceRepositoryMock_WithPagination_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int), args[1].(int))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithPagination_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithPagination_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithPagination_Call) RunAndReturn(run func(int, int) repo.Scope) *MasterServiceRepositoryMock_WithPagination_Call {
	_c.Call.Return(run)
	return _c
}

// WithSort provides a mock function with given fields: field, asc
func (_m *MasterServiceRepositoryMock) WithSort(field string, asc bool) repo.Scope {
	ret := _m.Called(field, asc)

	if len(ret) == 0 {
		panic("no return value specified for WithSort")
	}

	var r0 repo.Scope
	if rf, ok := ret.Get(0).(func(string, bool) repo.Scope); ok {
		r0 = rf(field, asc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(repo.Scope)
		}
	}

	return r0
}

// MasterServiceRepositoryMock_WithSort_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithSort'
type MasterServiceRepositoryMock_WithSort_Call struct {
	*mock.Call
}

// WithSort is a helper method to define mock.On call
//   - field string
//   - asc bool
func (_e *MasterServiceRepositoryMock_Expecter) WithSort(field interface{}, asc interface{}) *MasterServiceRepositoryMock_WithSort_Call {
	return &MasterServiceRepositoryMock_WithSort_Call{Call: _e.mock.On("WithSort", field, asc)}
}

func (_c *MasterServiceRepositoryMock_WithSort_Call) Run(run func(field string, asc bool)) *MasterServiceRepositoryMock_WithSort_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(bool))
	})
	return _c
}

func (_c *MasterServiceRepositoryMock_WithSort_Call) Return(_a0 repo.Scope) *MasterServiceRepositoryMock_WithSort_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MasterServiceRepositoryMock_WithSort_Call) RunAndReturn(run func(string, bool) repo.Scope) *MasterServiceRepositoryMock_WithSort_Call {
	_c.Call.Return(run)
	return _c
}

// NewMasterServiceRepositoryMock creates a new instance of MasterServiceRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMasterServiceRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *MasterServiceRepositoryMock {
	mock := &MasterServiceRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
