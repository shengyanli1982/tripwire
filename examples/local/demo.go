package main

import (
	"errors"
	"fmt"
	"time"

	tp "github.com/shengyanli1982/tripwire"
	cb "github.com/shengyanli1982/tripwire/circuitbreaker"
)

// 定义两种错误类型
// Define two types of errors
var (
	execError = errors.New("execution error") // 执行错误
	fbError   = errors.New("fallback error")  // 回退错误
)

// demoCallback 结构体用于实现回调接口
// demoCallback struct is used to implement the callback interface
type demoCallback struct{}

// OnSuccess 打印成功的消息
// OnSuccess prints the success message
func (d *demoCallback) OnSuccess(opterr error) {
	fmt.Printf("OnSuccess: %v\n", opterr) // 打印成功消息
}

// OnFailure 打印失败的消息和原因
// OnFailure prints the failure message and reason
func (d *demoCallback) OnFailure(opterr, reason error) {
	fmt.Printf("OnFailure: %v, %v\n", opterr, reason) // 打印失败消息和原因
}

// OnAccept 方法在接受时被调用，打印接受的原因和熔断器比例，以及失败比例。
// The OnAccept method is called when accepted, printing the accepted reason, fuse ratio, and failure ratio.
func (d *demoCallback) OnAccept(reason error, fuse, failure float64) {
	fmt.Printf("OnAccept: %v, fuse ratio: %v, failure ratio %v\n", reason, fuse, failure)
}

func main() {
	// 创建新的熔断器配置和熔断器
	// Create new circuit breaker configuration and circuit breaker
	config := cb.NewConfig().WithCallback(&demoCallback{})                     // 创建新的熔断器配置
	breaker := tp.New(tp.NewConfig().WithBreaker(cb.NewGoogleBreaker(config))) // 创建新的熔断器
	defer breaker.Stop()                                                       // 确保熔断器在主函数结束时停止

	// 模拟运行10次，成功
	// Simulate running 10 times, success
	for i := 0; i < 10; i++ {
		_ = breaker.Do(func() error {
			return nil // 返回nil表示成功
		})
	}

	// 案例1：默认情况下成功执行。
	// Case 1: Successful execution with default.
	fn := func() error {
		return nil // 返回nil表示成功
	}
	err := breaker.Do(fn) // 执行函数
	if err != nil {
		fmt.Printf("#Case1: Unexpected error: %v\n", err) // 如果有错误，打印错误
	} else {
		fmt.Printf("#Case1: Successful execution with default.\n") // 如果没有错误，打印成功消息
	}

	// 案例2：默认情况下执行失败。
	// Case 2: Failed execution with default.
	fn = func() error {
		return execError // 返回执行错误
	}
	err = breaker.Do(fn) // 执行函数
	if err != nil {
		fmt.Printf("#Case2: Unexpected error: %v\n", err) // 如果有错误，打印错误
	}

	// 案例3：执行失败，错误不可接受。
	// Case 3: Failed execution with unacceptable.
	acceptable := func(err error) bool {
		return !errors.Is(err, execError) // 如果错误是执行错误，返回false，否则返回true
	}
	err = breaker.DoWithAcceptable(fn, acceptable) // 执行函数，使用自定义的错误接受函数
	if err != nil {
		fmt.Printf("#Case3: Unexpected error: %v\n", err) // 如果有错误，打印错误
	} else {
		fmt.Printf("#Case3: Failed execution with unacceptable.\n") // 如果没有错误，打印执行失败的消息
	}

	// 案例4：执行失败，错误可接受。
	// Case 4: Failed execution with acceptable.
	acceptable = func(err error) bool {
		return errors.Is(err, execError) // 如果错误是执行错误，返回true，否则返回false
	}
	err = breaker.DoWithAcceptable(fn, acceptable) // 执行函数，使用自定义的错误接受函数
	if err != nil {
		fmt.Printf("#Case4: Unexpected error: %v\n", err) // 如果有错误，打印错误
	} else {
		fmt.Printf("#Case4: Failed execution with acceptable.\n") // 如果没有错误，打印执行失败的消息
	}

	// 模拟运行20次，失败
	// Simulate running 20 times, failed
	for i := 0; i < 20; i++ {
		_ = breaker.Do(func() error {
			return execError // 返回执行错误
		})
	}

	// 案例5：执行失败，有回退函数。
	// Case 5: Failed execution with fallback.
	fallback := func(err error) error {
		return fbError // 返回回退错误
	}
	err = breaker.DoWithFallback(fn, fallback) // 执行函数，使用自定义的回退函数
	if err != nil {
		fmt.Printf("#Case5: Unexpected error: %v\n", err) // 如果有错误，打印错误
	} else {
		fmt.Printf("#Case5: Failed execution with fallback.\n") // 如果没有错误，打印执行失败的消息
	}

	// 案例6：空闲5秒，成功执行。
	// Case 6: Idle for 5 seconds, successful execution.
	time.Sleep(5 * time.Second) // 等待5秒
	fn = func() error {
		return nil // 返回nil表示成功
	}
	err = breaker.Do(fn) // 执行函数
	if err != nil {
		fmt.Printf("#Case6: Unexpected error: %v\n", err) // 如果有错误，打印错误
	} else {
		fmt.Printf("#Case6: Idle for 5 seconds, successful execution.\n") // 如果没有错误，打印成功消息
	}
}
