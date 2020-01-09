// Code generated by mockery v1.0.0. DO NOT EDIT.

package securityscan

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SecretStore is an autogenerated mock type for the SecretStore type
type SecretStore struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, secretID
func (_m *SecretStore) Delete(ctx context.Context, secretID string) error {
	ret := _m.Called(ctx, secretID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, secretID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetIDByName provides a mock function with given fields: ctx, secretName
func (_m *SecretStore) GetIDByName(ctx context.Context, secretName string) (string, error) {
	ret := _m.Called(ctx, secretName)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, secretName)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, secretName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetNameByID provides a mock function with given fields: ctx, secretID
func (_m *SecretStore) GetNameByID(ctx context.Context, secretID string) (string, error) {
	ret := _m.Called(ctx, secretID)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, secretID)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, secretID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSecretValues provides a mock function with given fields: ctx, secretID
func (_m *SecretStore) GetSecretValues(ctx context.Context, secretID string) (map[string]string, error) {
	ret := _m.Called(ctx, secretID)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]string); ok {
		r0 = rf(ctx, secretID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, secretID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}