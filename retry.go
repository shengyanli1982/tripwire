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

// 定义空重试结构体
// Define the emptyRetry struct
type emptyRetry struct{}

// TryOnConflict 方法尝试执行给定的函数，并返回结果
// The TryOnConflict method tries to execute the given function and returns the result
func (r *emptyRetry) TryOnConflict(fn com.RetryableFunc) com.RetryResult {
	re := result{count: 1}
	re.data, re.tryError = fn()
	return &re
}

// NewEmptyRetry 函数创建并返回一个新的空重试实例
// The NewEmptyRetry function creates and returns a new instance of emptyRetry
func NewEmptyRetry() com.Retry {
	return &emptyRetry{}
}
