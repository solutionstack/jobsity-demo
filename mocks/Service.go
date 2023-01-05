// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	models "github.com/solutionstack/jobsity-demo/models"
	mock "github.com/stretchr/testify/mock"

	testing "testing"

	uuid "github.com/google/uuid"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: user
func (_m *Service) CreateUser(user models.Signup) (uuid.UUID, error) {
	ret := _m.Called(user)

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func(models.Signup) uuid.UUID); ok {
		r0 = rf(user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Signup) error); ok {
		r1 = rf(user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateLogin provides a mock function with given fields: login
func (_m *Service) ValidateLogin(login models.Login) (*models.UserRecord, error) {
	ret := _m.Called(login)

	var r0 *models.UserRecord
	if rf, ok := ret.Get(0).(func(models.Login) *models.UserRecord); ok {
		r0 = rf(login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(models.Login) error); ok {
		r1 = rf(login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewService creates a new instance of Service. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t testing.TB) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}