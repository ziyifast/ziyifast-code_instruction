package main

import (
	"errors"
	"fmt"
)

/*
	List:
		- NewList(): 创建一个新的列表
		- Add(element): 在列表末尾添加元素
		- Remove(index): 根据索引移除元素
		- Size(): 返回列表的大小
		- Get(index): 根据索引获取元素
		- IsEmpty(): 判断列表是否为空
		- Clear(): 清空列表
		- GetFirst(): 获取第一个元素
		- GetLast(): 获取最后一个元素
		- RemoveLast(): 移除最后一个元素
		- AddFirst(element): 在列表开头添加元素
		- RemoveFirst(): 移除第一个元素
*/

type List struct {
	data []interface{}
}

// NewList 创建一个新的列表
func NewList() *List {
	return &List{
		data: []interface{}{},
	}
}

// Add 在列表末尾添加元素
func (l *List) Add(v interface{}) {
	l.data = append(l.data, v)
}

// Remove 根据索引移除元素
func (l *List) Remove(index int) error {
	if index < 0 || index >= len(l.data) {
		return errors.New("index out of bounds")
	}
	l.data = append(l.data[:index], l.data[index+1:]...)
	return nil
}

// Size 返回列表的大小
func (l *List) Size() int {
	return len(l.data)
}

// Get 根据索引获取元素
func (l *List) Get(index int) (interface{}, error) {
	if index < 0 || index >= len(l.data) {
		return nil, errors.New("index out of bounds")
	}
	return l.data[index], nil
}

// IsEmpty 判断列表是否为空
func (l *List) IsEmpty() bool {
	return len(l.data) == 0
}

// Clear 清空列表
func (l *List) Clear() {
	l.data = []interface{}{}
}

// GetFirst 获取第一个元素
func (l *List) GetFirst() (interface{}, error) {
	if l.IsEmpty() {
		return nil, errors.New("list is empty")
	}
	return l.data[0], nil
}

// GetLast 获取最后一个元素
func (l *List) GetLast() (interface{}, error) {
	if l.IsEmpty() {
		return nil, errors.New("list is empty")
	}
	return l.data[len(l.data)-1], nil
}

// AddFirst 在列表开头添加元素
func (l *List) AddFirst(v interface{}) {
	l.data = append([]interface{}{v}, l.data...)
}

// RemoveFirst 移除第一个元素
func (l *List) RemoveFirst() error {
	if l.IsEmpty() {
		return errors.New("list is empty")
	}
	l.data = l.data[1:]
	return nil
}

// RemoveLast 移除最后一个元素
func (l *List) RemoveLast() error {
	if l.IsEmpty() {
		return errors.New("list is empty")
	}
	l.data = l.data[:len(l.data)-1]
	return nil
}

func main() {
	list := NewList()

	// 测试 Add 和 Get
	list.Add(1)
	list.Add(2)
	list.Add(3)
	value, err := list.Get(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Value at index 1:", value) // 输出: Value at index 1: 2
	}

	// 测试 Remove
	err = list.Remove(1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Size after remove:", list.Size()) // 输出: Size after remove: 2
	}

	// 测试 GetFirst 和 GetLast
	first, err := list.GetFirst()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("First element:", first) // 输出: First element: 1
	}

	last, err := list.GetLast()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Last element:", last) // 输出: Last element: 3
	}

	// 测试 AddFirst 和 RemoveFirst
	list.AddFirst(0)
	first, err = list.GetFirst()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("First element after addFirst:", first) // 输出: First element after addFirst: 0
	}

	err = list.RemoveFirst()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Size after removeFirst:", list.Size()) // 输出: Size after removeFirst: 2
	}

	// 测试 RemoveLast
	err = list.RemoveLast()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Size after removeLast:", list.Size()) // 输出: Size after removeLast: 1
	}

	// 测试 Clear
	list.Clear()
	fmt.Println("Is list empty after clear?", list.IsEmpty()) // 输出: Is list empty after clear? true
}
