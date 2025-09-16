// Package config
package config

type validType int

const (
	PASS validType = iota
	FAIL
)

type ValidResult struct {
	validType validType
	err       error
	originErr error
}

func ValidPass() *ValidResult {
	return &ValidResult{validType: PASS, err: nil, originErr: nil}
}

func ValidFail(err error) *ValidResult {
	return &ValidResult{validType: FAIL, err: err, originErr: nil}
}

func ValidFailWith(err error, originErr error) *ValidResult {
	return &ValidResult{validType: FAIL, err: err, originErr: originErr}
}

func (r *ValidResult) IsFail() bool {
	return r.validType == FAIL
}

func (r *ValidResult) Err() error {
	return r.err
}

func (r *ValidResult) Error() string {
	return r.err.Error()
}

func (r *ValidResult) OriginErr() error { return r.originErr }
