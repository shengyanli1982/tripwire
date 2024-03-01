package tripwire

import (
	com "github.com/shengyanli1982/tripwire/common"
)

// 定义结果结构体
// Define the result struct
type result struct {
	count    uint64 // 计数器 Counter
	data     any    // 数据 Data
	tryError error  // 错误信息 Error information
	errs     []error
}

// Data 方法返回结果中的数据
// The Data method returns the data in the result
func (r *result) Data() any {
	return r.data
}

// TryError 方法返回结果中的错误信息
// The TryError method returns the error information in the result
func (r *result) TryError() error {
	return r.tryError
}

// IsSuccess 方法检查结果是否成功，如果没有错误，返回 true
// The IsSuccess method checks whether the result is successful, returns true if there is no error
func (r *result) IsSuccess() bool {
	return r.tryError == nil
}

// Count 方法返回结果的计数
// The Count method returns the count of the result
func (r *result) Count() int64 {
	return int64(r.count)
}

// ExecErrors 方法返回结果中的所有错误
// The ExecErrors method returns all errors in the result
func (r *result) ExecErrors() []error {
	return r.errs
}

// LastExecError 方法返回结果中的最后一个错误
// The LastExecError method returns the last error in the result
func (r *result) LastExecError() error {
	if len(r.errs) > 0 {
		return r.errs[len(r.errs)-1]
	}
	return nil
}

// FirstExecError 方法返回结果中的第一个错误
// The FirstExecError method returns the first error in the result
func (r *result) FirstExecError() error {
	if len(r.errs) > 0 {
		return r.errs[0]
	}
	return nil
}

// ExecErrorByIndex 方法返回结果中指定索引处的错误
// The ExecErrorByIndex method returns the error at the specified index in the result
func (r *result) ExecErrorByIndex(idx int) error {
	if idx >= 0 && idx < len(r.errs) {
		return r.errs[idx]
	}
	return nil
}

// 定义空重试结构体
// Define the emptyRetry struct
type emptyRetry struct{}

// TryOnConflictVal 方法尝试执行给定的函数，并返回结果
// The TryOnConflictVal method tries to execute the given function and returns the result
func (r *emptyRetry) TryOnConflictVal(fn com.RetryableFunc) com.RetryResult {
	re := result{count: 1}
	re.data, re.tryError = fn()
	return &re
}

// NewEmptyRetry 函数创建并返回一个新的空重试实例
// The NewEmptyRetry function creates and returns a new instance of emptyRetry
func NewEmptyRetry() com.Retry {
	return &emptyRetry{}
}
