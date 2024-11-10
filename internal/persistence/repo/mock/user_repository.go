// Code generated by mockery v2.46.3. DO NOT EDIT.

package mock

import (
	context "context"

	model "github.com/mandarine-io/Backend/internal/persistence/model"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// UserRepositoryMock is an autogenerated mock type for the UserRepository type
type UserRepositoryMock struct {
	mock.Mock
}

type UserRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *UserRepositoryMock) EXPECT() *UserRepositoryMock_Expecter {
	return &UserRepositoryMock_Expecter{mock: &_m.Mock}
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *UserRepositoryMock) CreateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntity) (*model.UserEntity, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntity) *model.UserEntity); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserEntity) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_CreateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateUser'
type UserRepositoryMock_CreateUser_Call struct {
	*mock.Call
}

// CreateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user *model.UserEntity
func (_e *UserRepositoryMock_Expecter) CreateUser(ctx interface{}, user interface{}) *UserRepositoryMock_CreateUser_Call {
	return &UserRepositoryMock_CreateUser_Call{Call: _e.mock.On("CreateUser", ctx, user)}
}

func (_c *UserRepositoryMock_CreateUser_Call) Run(run func(ctx context.Context, user *model.UserEntity)) *UserRepositoryMock_CreateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserEntity))
	})
	return _c
}

func (_c *UserRepositoryMock_CreateUser_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_CreateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_CreateUser_Call) RunAndReturn(run func(context.Context, *model.UserEntity) (*model.UserEntity, error)) *UserRepositoryMock_CreateUser_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteExpiredUser provides a mock function with given fields: ctx
func (_m *UserRepositoryMock) DeleteExpiredUser(ctx context.Context) (*model.UserEntity, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExpiredUser")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*model.UserEntity, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *model.UserEntity); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_DeleteExpiredUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteExpiredUser'
type UserRepositoryMock_DeleteExpiredUser_Call struct {
	*mock.Call
}

// DeleteExpiredUser is a helper method to define mock.On call
//   - ctx context.Context
func (_e *UserRepositoryMock_Expecter) DeleteExpiredUser(ctx interface{}) *UserRepositoryMock_DeleteExpiredUser_Call {
	return &UserRepositoryMock_DeleteExpiredUser_Call{Call: _e.mock.On("DeleteExpiredUser", ctx)}
}

func (_c *UserRepositoryMock_DeleteExpiredUser_Call) Run(run func(ctx context.Context)) *UserRepositoryMock_DeleteExpiredUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *UserRepositoryMock_DeleteExpiredUser_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_DeleteExpiredUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_DeleteExpiredUser_Call) RunAndReturn(run func(context.Context) (*model.UserEntity, error)) *UserRepositoryMock_DeleteExpiredUser_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsUserByEmail provides a mock function with given fields: ctx, email
func (_m *UserRepositoryMock) ExistsUserByEmail(ctx context.Context, email string) (bool, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for ExistsUserByEmail")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_ExistsUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsUserByEmail'
type UserRepositoryMock_ExistsUserByEmail_Call struct {
	*mock.Call
}

// ExistsUserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *UserRepositoryMock_Expecter) ExistsUserByEmail(ctx interface{}, email interface{}) *UserRepositoryMock_ExistsUserByEmail_Call {
	return &UserRepositoryMock_ExistsUserByEmail_Call{Call: _e.mock.On("ExistsUserByEmail", ctx, email)}
}

func (_c *UserRepositoryMock_ExistsUserByEmail_Call) Run(run func(ctx context.Context, email string)) *UserRepositoryMock_ExistsUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByEmail_Call) Return(_a0 bool, _a1 error) *UserRepositoryMock_ExistsUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByEmail_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *UserRepositoryMock_ExistsUserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsUserById provides a mock function with given fields: ctx, id
func (_m *UserRepositoryMock) ExistsUserById(ctx context.Context, id uuid.UUID) (bool, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for ExistsUserById")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (bool, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_ExistsUserById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsUserById'
type UserRepositoryMock_ExistsUserById_Call struct {
	*mock.Call
}

// ExistsUserById is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
func (_e *UserRepositoryMock_Expecter) ExistsUserById(ctx interface{}, id interface{}) *UserRepositoryMock_ExistsUserById_Call {
	return &UserRepositoryMock_ExistsUserById_Call{Call: _e.mock.On("ExistsUserById", ctx, id)}
}

func (_c *UserRepositoryMock_ExistsUserById_Call) Run(run func(ctx context.Context, id uuid.UUID)) *UserRepositoryMock_ExistsUserById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *UserRepositoryMock_ExistsUserById_Call) Return(_a0 bool, _a1 error) *UserRepositoryMock_ExistsUserById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_ExistsUserById_Call) RunAndReturn(run func(context.Context, uuid.UUID) (bool, error)) *UserRepositoryMock_ExistsUserById_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsUserByUsername provides a mock function with given fields: ctx, username
func (_m *UserRepositoryMock) ExistsUserByUsername(ctx context.Context, username string) (bool, error) {
	ret := _m.Called(ctx, username)

	if len(ret) == 0 {
		panic("no return value specified for ExistsUserByUsername")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, username)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_ExistsUserByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsUserByUsername'
type UserRepositoryMock_ExistsUserByUsername_Call struct {
	*mock.Call
}

// ExistsUserByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
func (_e *UserRepositoryMock_Expecter) ExistsUserByUsername(ctx interface{}, username interface{}) *UserRepositoryMock_ExistsUserByUsername_Call {
	return &UserRepositoryMock_ExistsUserByUsername_Call{Call: _e.mock.On("ExistsUserByUsername", ctx, username)}
}

func (_c *UserRepositoryMock_ExistsUserByUsername_Call) Run(run func(ctx context.Context, username string)) *UserRepositoryMock_ExistsUserByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByUsername_Call) Return(_a0 bool, _a1 error) *UserRepositoryMock_ExistsUserByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByUsername_Call) RunAndReturn(run func(context.Context, string) (bool, error)) *UserRepositoryMock_ExistsUserByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// ExistsUserByUsernameOrEmail provides a mock function with given fields: ctx, username, email
func (_m *UserRepositoryMock) ExistsUserByUsernameOrEmail(ctx context.Context, username string, email string) (bool, error) {
	ret := _m.Called(ctx, username, email)

	if len(ret) == 0 {
		panic("no return value specified for ExistsUserByUsernameOrEmail")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (bool, error)); ok {
		return rf(ctx, username, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, username, email)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, username, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_ExistsUserByUsernameOrEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ExistsUserByUsernameOrEmail'
type UserRepositoryMock_ExistsUserByUsernameOrEmail_Call struct {
	*mock.Call
}

// ExistsUserByUsernameOrEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - email string
func (_e *UserRepositoryMock_Expecter) ExistsUserByUsernameOrEmail(ctx interface{}, username interface{}, email interface{}) *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call {
	return &UserRepositoryMock_ExistsUserByUsernameOrEmail_Call{Call: _e.mock.On("ExistsUserByUsernameOrEmail", ctx, username, email)}
}

func (_c *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call) Run(run func(ctx context.Context, username string, email string)) *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call) Return(_a0 bool, _a1 error) *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call) RunAndReturn(run func(context.Context, string, string) (bool, error)) *UserRepositoryMock_ExistsUserByUsernameOrEmail_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserByEmail provides a mock function with given fields: ctx, email, rolePreload
func (_m *UserRepositoryMock) FindUserByEmail(ctx context.Context, email string, rolePreload bool) (*model.UserEntity, error) {
	ret := _m.Called(ctx, email, rolePreload)

	if len(ret) == 0 {
		panic("no return value specified for FindUserByEmail")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (*model.UserEntity, error)); ok {
		return rf(ctx, email, rolePreload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) *model.UserEntity); ok {
		r0 = rf(ctx, email, rolePreload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, email, rolePreload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_FindUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserByEmail'
type UserRepositoryMock_FindUserByEmail_Call struct {
	*mock.Call
}

// FindUserByEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
//   - rolePreload bool
func (_e *UserRepositoryMock_Expecter) FindUserByEmail(ctx interface{}, email interface{}, rolePreload interface{}) *UserRepositoryMock_FindUserByEmail_Call {
	return &UserRepositoryMock_FindUserByEmail_Call{Call: _e.mock.On("FindUserByEmail", ctx, email, rolePreload)}
}

func (_c *UserRepositoryMock_FindUserByEmail_Call) Run(run func(ctx context.Context, email string, rolePreload bool)) *UserRepositoryMock_FindUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *UserRepositoryMock_FindUserByEmail_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_FindUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_FindUserByEmail_Call) RunAndReturn(run func(context.Context, string, bool) (*model.UserEntity, error)) *UserRepositoryMock_FindUserByEmail_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserById provides a mock function with given fields: ctx, id, rolePreload
func (_m *UserRepositoryMock) FindUserById(ctx context.Context, id uuid.UUID, rolePreload bool) (*model.UserEntity, error) {
	ret := _m.Called(ctx, id, rolePreload)

	if len(ret) == 0 {
		panic("no return value specified for FindUserById")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) (*model.UserEntity, error)); ok {
		return rf(ctx, id, rolePreload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, bool) *model.UserEntity); ok {
		r0 = rf(ctx, id, rolePreload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, bool) error); ok {
		r1 = rf(ctx, id, rolePreload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_FindUserById_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserById'
type UserRepositoryMock_FindUserById_Call struct {
	*mock.Call
}

// FindUserById is a helper method to define mock.On call
//   - ctx context.Context
//   - id uuid.UUID
//   - rolePreload bool
func (_e *UserRepositoryMock_Expecter) FindUserById(ctx interface{}, id interface{}, rolePreload interface{}) *UserRepositoryMock_FindUserById_Call {
	return &UserRepositoryMock_FindUserById_Call{Call: _e.mock.On("FindUserById", ctx, id, rolePreload)}
}

func (_c *UserRepositoryMock_FindUserById_Call) Run(run func(ctx context.Context, id uuid.UUID, rolePreload bool)) *UserRepositoryMock_FindUserById_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(bool))
	})
	return _c
}

func (_c *UserRepositoryMock_FindUserById_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_FindUserById_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_FindUserById_Call) RunAndReturn(run func(context.Context, uuid.UUID, bool) (*model.UserEntity, error)) *UserRepositoryMock_FindUserById_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserByUsername provides a mock function with given fields: ctx, username, rolePreload
func (_m *UserRepositoryMock) FindUserByUsername(ctx context.Context, username string, rolePreload bool) (*model.UserEntity, error) {
	ret := _m.Called(ctx, username, rolePreload)

	if len(ret) == 0 {
		panic("no return value specified for FindUserByUsername")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (*model.UserEntity, error)); ok {
		return rf(ctx, username, rolePreload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) *model.UserEntity); ok {
		r0 = rf(ctx, username, rolePreload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, username, rolePreload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_FindUserByUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserByUsername'
type UserRepositoryMock_FindUserByUsername_Call struct {
	*mock.Call
}

// FindUserByUsername is a helper method to define mock.On call
//   - ctx context.Context
//   - username string
//   - rolePreload bool
func (_e *UserRepositoryMock_Expecter) FindUserByUsername(ctx interface{}, username interface{}, rolePreload interface{}) *UserRepositoryMock_FindUserByUsername_Call {
	return &UserRepositoryMock_FindUserByUsername_Call{Call: _e.mock.On("FindUserByUsername", ctx, username, rolePreload)}
}

func (_c *UserRepositoryMock_FindUserByUsername_Call) Run(run func(ctx context.Context, username string, rolePreload bool)) *UserRepositoryMock_FindUserByUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *UserRepositoryMock_FindUserByUsername_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_FindUserByUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_FindUserByUsername_Call) RunAndReturn(run func(context.Context, string, bool) (*model.UserEntity, error)) *UserRepositoryMock_FindUserByUsername_Call {
	_c.Call.Return(run)
	return _c
}

// FindUserByUsernameOrEmail provides a mock function with given fields: ctx, login, rolePreload
func (_m *UserRepositoryMock) FindUserByUsernameOrEmail(ctx context.Context, login string, rolePreload bool) (*model.UserEntity, error) {
	ret := _m.Called(ctx, login, rolePreload)

	if len(ret) == 0 {
		panic("no return value specified for FindUserByUsernameOrEmail")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) (*model.UserEntity, error)); ok {
		return rf(ctx, login, rolePreload)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, bool) *model.UserEntity); ok {
		r0 = rf(ctx, login, rolePreload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, bool) error); ok {
		r1 = rf(ctx, login, rolePreload)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_FindUserByUsernameOrEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindUserByUsernameOrEmail'
type UserRepositoryMock_FindUserByUsernameOrEmail_Call struct {
	*mock.Call
}

// FindUserByUsernameOrEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - login string
//   - rolePreload bool
func (_e *UserRepositoryMock_Expecter) FindUserByUsernameOrEmail(ctx interface{}, login interface{}, rolePreload interface{}) *UserRepositoryMock_FindUserByUsernameOrEmail_Call {
	return &UserRepositoryMock_FindUserByUsernameOrEmail_Call{Call: _e.mock.On("FindUserByUsernameOrEmail", ctx, login, rolePreload)}
}

func (_c *UserRepositoryMock_FindUserByUsernameOrEmail_Call) Run(run func(ctx context.Context, login string, rolePreload bool)) *UserRepositoryMock_FindUserByUsernameOrEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(bool))
	})
	return _c
}

func (_c *UserRepositoryMock_FindUserByUsernameOrEmail_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_FindUserByUsernameOrEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_FindUserByUsernameOrEmail_Call) RunAndReturn(run func(context.Context, string, bool) (*model.UserEntity, error)) *UserRepositoryMock_FindUserByUsernameOrEmail_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUser provides a mock function with given fields: ctx, user
func (_m *UserRepositoryMock) UpdateUser(ctx context.Context, user *model.UserEntity) (*model.UserEntity, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *model.UserEntity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntity) (*model.UserEntity, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UserEntity) *model.UserEntity); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.UserEntity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UserEntity) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserRepositoryMock_UpdateUser_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUser'
type UserRepositoryMock_UpdateUser_Call struct {
	*mock.Call
}

// UpdateUser is a helper method to define mock.On call
//   - ctx context.Context
//   - user *model.UserEntity
func (_e *UserRepositoryMock_Expecter) UpdateUser(ctx interface{}, user interface{}) *UserRepositoryMock_UpdateUser_Call {
	return &UserRepositoryMock_UpdateUser_Call{Call: _e.mock.On("UpdateUser", ctx, user)}
}

func (_c *UserRepositoryMock_UpdateUser_Call) Run(run func(ctx context.Context, user *model.UserEntity)) *UserRepositoryMock_UpdateUser_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*model.UserEntity))
	})
	return _c
}

func (_c *UserRepositoryMock_UpdateUser_Call) Return(_a0 *model.UserEntity, _a1 error) *UserRepositoryMock_UpdateUser_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *UserRepositoryMock_UpdateUser_Call) RunAndReturn(run func(context.Context, *model.UserEntity) (*model.UserEntity, error)) *UserRepositoryMock_UpdateUser_Call {
	_c.Call.Return(run)
	return _c
}

// NewUserRepositoryMock creates a new instance of UserRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepositoryMock {
	mock := &UserRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
