package tripwire

import (
	com "github.com/shengyanli1982/tripwire/common"
)

type result struct {
	count    uint64
	data     any
	tryError error
}

func (r *result) Data() any {
	return r.data
}

func (r *result) TryError() error {
	return r.tryError
}

func (r *result) IsSuccess() bool {
	return r.tryError == nil
}

func (r *result) Count() int64 {
	return int64(r.count)
}

type emptyRetry struct{}

func (r *emptyRetry) TryOnConflict(fn com.RetryableFunc) com.RetryResult {
	return &result{}
}

func NewEmptyRetry() com.Retry {
	return &emptyRetry{}
}
