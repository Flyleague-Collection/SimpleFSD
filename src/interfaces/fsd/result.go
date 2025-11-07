// Package fsd
package fsd

import "errors"

type Result struct {
	Success bool
	Errno   ClientError
	Fatal   bool
	Env     string
	Err     error
}

var defaultError = errors.New("no details error provided")

func ResultSuccess() *Result {
	return &Result{
		Success: true,
		Errno:   CommandOk,
		Fatal:   false,
		Env:     "",
		Err:     nil,
	}
}

func ResultError(errno ClientError, fatal bool, env string, err error) *Result {
	result := &Result{
		Success: false,
		Errno:   errno,
		Fatal:   fatal,
		Env:     env,
	}
	if err != nil {
		result.Err = err
	} else {
		result.Err = defaultError
	}
	return result
}
