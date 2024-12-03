package main

import (
	"fmt"
)

/*
  Stack:
	- Push(item): 入栈
	- Pop(): 出栈
	- Peek(): 返回栈顶元素，但不删除它
	- IsEmpty(): 判断栈是否为空
	- Search(item): 搜索 item 元素在栈中的位置，如果没找到，返回 -1
	- Clear(): 清空栈
*/

type Stack struct {
	data []interface{}
}

func NewStack() *Stack {
	return &Stack{
		data: []interface{}{},
	}
}

// Push 入栈
func (s *Stack) Push(v interface{}) {
	s.data = append(s.data, v)
}

// Pop 出栈
func (s *Stack) Pop() interface{} {
	if len(s.data) == 0 {
		return nil
	}
	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

// Peek 返回栈顶元素，但不删除它
func (s *Stack) Peek() interface{} {
	if len(s.data) == 0 {
		return nil
	}
	return s.data[len(s.data)-1]
}

// IsEmpty 判断栈是否为空
func (s *Stack) IsEmpty() bool {
	return len(s.data) == 0
}

// Search 搜索 item 元素在栈中的位置，如果没找到，返回 -1
func (s *Stack) Search(v interface{}) int {
	for index, value := range s.data {
		if value == v {
			return index
		}
	}
	return -1
}

// Clear 清空栈
func (s *Stack) Clear() {
	s.data = []interface{}{}
}

func main() {
	stack := NewStack()

	// 测试 Push 和 Peek
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	fmt.Println("Top element:", stack.Peek()) // 输出: Top element: 3

	// 测试 Pop
	fmt.Println("Popped element:", stack.Pop())         // 输出: Popped element: 3
	fmt.Println("Top element after pop:", stack.Peek()) // 输出: Top element after pop: 2

	// 测试 IsEmpty
	fmt.Println("Is stack empty?", stack.IsEmpty()) // 输出: Is stack empty? false

	// 测试 Search
	fmt.Println("Index of 2:", stack.Search(2)) // 输出: Index of 2: 1
	fmt.Println("Index of 3:", stack.Search(3)) // 输出: Index of 3: -1

	// 测试 Clear
	stack.Clear()
	fmt.Println("Is stack empty after clear?", stack.IsEmpty()) // 输出: Is stack empty after clear? true
}
