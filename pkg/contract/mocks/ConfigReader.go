// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	contract "github.com/lingwei0604/kitty/pkg/contract"

	time "time"
)

// ConfigReader is an autogenerated mock type for the ConfigReader type
type ConfigReader struct {
	mock.Mock
}

// Bool provides a mock function with given fields: _a0
func (_m *ConfigReader) Bool(_a0 string) bool {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Cut provides a mock function with given fields: _a0
func (_m *ConfigReader) Cut(_a0 string) contract.ConfigReader {
	ret := _m.Called(_a0)

	var r0 contract.ConfigReader
	if rf, ok := ret.Get(0).(func(string) contract.ConfigReader); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(contract.ConfigReader)
		}
	}

	return r0
}

// Duration provides a mock function with given fields: _a0
func (_m *ConfigReader) Duration(_a0 string) time.Duration {
	ret := _m.Called(_a0)

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func(string) time.Duration); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// Float64 provides a mock function with given fields: _a0
func (_m *ConfigReader) Float64(_a0 string) float64 {
	ret := _m.Called(_a0)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(float64)
	}

	return r0
}

// Get provides a mock function with given fields: _a0
func (_m *ConfigReader) Get(_a0 string) interface{} {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Int provides a mock function with given fields: _a0
func (_m *ConfigReader) Int(_a0 string) int {
	ret := _m.Called(_a0)

	var r0 int
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// String provides a mock function with given fields: _a0
func (_m *ConfigReader) String(_a0 string) string {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Strings provides a mock function with given fields: _a0
func (_m *ConfigReader) Strings(_a0 string) []string {
	ret := _m.Called(_a0)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// Unmarshal provides a mock function with given fields: path, o
func (_m *ConfigReader) Unmarshal(path string, o interface{}) error {
	ret := _m.Called(path, o)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(path, o)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
