// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package app

import (
	context "context"

	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	mock "github.com/stretchr/testify/mock"
)

// mockS3StorageClientAPI is an autogenerated mock type for the s3StorageClientAPI type
type mockS3StorageClientAPI struct {
	mock.Mock
}

// DeleteObject provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockS3StorageClientAPI) DeleteObject(_a0 context.Context, _a1 *s3.DeleteObjectInput, _a2 ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.DeleteObjectOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) *s3.DeleteObjectOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.DeleteObjectOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetObject provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockS3StorageClientAPI) GetObject(_a0 context.Context, _a1 *s3.GetObjectInput, _a2 ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.GetObjectOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) *s3.GetObjectOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.GetObjectOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HeadObject provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockS3StorageClientAPI) HeadObject(_a0 context.Context, _a1 *s3.HeadObjectInput, _a2 ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.HeadObjectOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) *s3.HeadObjectOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.HeadObjectOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutObject provides a mock function with given fields: _a0, _a1, _a2
func (_m *mockS3StorageClientAPI) PutObject(_a0 context.Context, _a1 *s3.PutObjectInput, _a2 ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *s3.PutObjectOutput
	if rf, ok := ret.Get(0).(func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) *s3.PutObjectOutput); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s3.PutObjectOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
