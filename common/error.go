package common

import "errors"

var (
	// 定义服务不可用的错误
	// Define the error for service unavailable
	ErrorServiceUnavailable = errors.New("service unavailable")

	// 滚动窗口停止的错误。
	// Error when the rolling window is stopped.
	ErrorRollingWindowStopped = errors.New("rolling window stopped")
)
