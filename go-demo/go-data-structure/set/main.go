package main

import (
	"fmt"
	"sync"
)

/*
 Set: 可以去除重复元素
	- Add: 添加元素
	- Remove: 删除元素
	- Contains: 检查元素是否存在
	- IsEmpty: 判断集合是否为空
	- Clear: 清空集合
	- Iterator: 返回一个迭代器通道
*/

type Set struct {
	mu   sync.RWMutex
	data map[interface{}]bool
}

// NewSet 创建一个新的集合
func NewSet() *Set {
	return &Set{
		data: make(map[interface{}]bool),
	}
}

// Add 添加元素到集合
func (s *Set) Add(value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[value] = true
}

// Remove 从集合中删除元素
func (s *Set) Remove(value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, value)
}

// Contains 检查元素是否存在于集合中
func (s *Set) Contains(value interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data[value]
}

// IsEmpty 判断集合是否为空
func (s *Set) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data) == 0
}

// Clear 清空集合
func (s *Set) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[interface{}]bool)
}

// Iterator 返回一个迭代器通道
func (s *Set) Iterator() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in Iterator:", r)
			}
			close(ch)
		}()
		s.mu.RLock()
		defer s.mu.RUnlock()
		for k := range s.data {
			ch <- k
		}
	}()
	return ch
}

func main() {
	set := NewSet()

	// 测试 Add 和 Contains
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println("Contains 1:", set.Contains(1)) // 输出: Contains 1: true
	fmt.Println("Contains 4:", set.Contains(4)) // 输出: Contains 4: false

	// 测试 Remove
	set.Remove(2)
	fmt.Println("Contains 2 after remove:", set.Contains(2)) // 输出: Contains 2 after remove: false

	// 测试 IsEmpty
	fmt.Println("Is set empty?", set.IsEmpty()) // 输出: Is set empty? false

	// 测试 Clear
	set.Clear()
	fmt.Println("Is set empty after clear?", set.IsEmpty()) // 输出: Is set empty after clear? true

	// 测试 Iterator
	set.Add(1)
	set.Add(2)
	set.Add(3)
	fmt.Println("Elements in set:")
	for i := range set.Iterator() {
		fmt.Println(i)
	}
	// 其他测试代码
	data := make([]int, 2, 20)
	data[0] = -1
	fmt.Println("Length of data:", len(data))   // 输出: Length of data: 2
	fmt.Println("Capacity of data:", cap(data)) // 输出: Capacity of data: 20
}
