// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Store is an autogenerated mock type for the store type
type Store struct {
	mock.Mock
}

type Store_Expecter struct {
	mock *mock.Mock
}

func (_m *Store) EXPECT() *Store_Expecter {
	return &Store_Expecter{mock: &_m.Mock}
}

// Check provides a mock function with given fields:
func (_m *Store) Check() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Check")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store_Check_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Check'
type Store_Check_Call struct {
	*mock.Call
}

// Check is a helper method to define mock.On call
func (_e *Store_Expecter) Check() *Store_Check_Call {
	return &Store_Check_Call{Call: _e.mock.On("Check")}
}

func (_c *Store_Check_Call) Run(run func()) *Store_Check_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Store_Check_Call) Return(_a0 error) *Store_Check_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Store_Check_Call) RunAndReturn(run func() error) *Store_Check_Call {
	_c.Call.Return(run)
	return _c
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
