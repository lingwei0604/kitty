// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// FileRepository is an autogenerated mock type for the FileRepository type
type FileRepository struct {
	mock.Mock
}

// UploadFromUrl provides a mock function with given fields: ctx, url
func (_m *FileRepository) UploadFromUrl(ctx context.Context, url string) (string, error) {
	ret := _m.Called(ctx, url)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, url)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, url)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
