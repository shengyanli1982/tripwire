// 定义包名
// Define the package name
package rolling

// 导入依赖
// Import dependencies
import "github.com/shengyanli1982/tripwire/internal/utils"

// 定义 Ring 结构体
// Define the Ring struct
type Ring struct {
	slots  []any  // 存储元素的数组 Array to store elements
	head   uint64 // 头部索引 Head index
	tail   uint64 // 尾部索引 Tail index
	count  int    // 当前元素数量 Current number of elements
	cap    int    // Ring 的容量 Capacity of the Ring
	optval uint64 // 用于优化的值 Value for optimization
}

// NewRing 函数创建并返回一个新的 Ring 实例
// The NewRing function creates and returns a new instance of Ring
func NewRing(cap int) *Ring {
	cap = utils.FindNextPowerOfTwo(cap)
	return &Ring{
		slots:  make([]any, cap),
		cap:    cap,
		optval: uint64(cap - 1),
	}
}

// Reset 方法重置 Ring 的 head、tail 和 count
// The Reset method resets the head, tail and count of the Ring
func (r *Ring) Reset() {
	r.head = 0
	r.tail = 0
	r.count = 0
}

// Len 方法返回 Ring 的长度，如果 count 小于 cap，则返回 count，否则返回 cap
// The Len method returns the length of the Ring, returns count if count is less than cap, otherwise returns cap
func (r *Ring) Len() int {
	if r.count < r.cap {
		return r.count
	}
	return r.cap
}

// Cap 方法返回 Ring 的容量
// The Cap method returns the capacity of the Ring
func (r *Ring) Cap() int {
	return r.cap
}

// Head 方法返回 Ring 的头部元素，如果 count 为 0，则返回 nil
// The Head method returns the head element of the Ring, returns nil if count is 0
func (r *Ring) Head() any {
	if r.count == 0 {
		return nil
	}
	return r.slots[r.head]
}

// Tail 方法返回 Ring 的尾部元素，如果 count 为 0，则返回 nil
// The Tail method returns the tail element of the Ring, returns nil if count is 0
func (r *Ring) Tail() any {
	if r.count == 0 {
		return nil
	}
	return r.slots[(r.tail-1)&r.optval]
}

// At 方法返回 Ring 中索引为 i 的元素，如果 count 为 0，则返回 nil
// The At method returns the element at index i in the Ring, returns nil if count is 0
func (r *Ring) At(i int) any {
	if r.count == 0 {
		return nil
	}
	return r.slots[(r.head+uint64(i))&r.optval]
}

// Push 方法将一个元素添加到 Ring 的尾部
// The Push method adds an element to the tail of the Ring
func (r *Ring) Push(v any) {
	r.slots[r.tail] = v
	r.tail = (r.tail + 1) & r.optval
	if r.count < r.cap {
		r.count++
	} else {
		r.head = (r.head + 1) & r.optval
	}
}

// Values 方法返回 Ring 中的所有元素
// The Values method returns all elements in the Ring
func (r *Ring) Values() []any {
	return r.slots
}
