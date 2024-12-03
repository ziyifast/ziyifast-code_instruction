package main

import (
	"container/list"
	"fmt"
)

/*
	Deque:
		- PushFront: 在队列前端插入元素
		- PushBack: 在队列后端插入元素
		- PopFront: 从队列前端移除并返回元素
		- PopBack: 从队列后端移除并返回元素
		...
*/

// Deque 双端队列结构体
type Deque struct {
	data *list.List
}

// NewDeque 创建一个新的双端队列
func NewDeque() *Deque {
	return &Deque{data: list.New()}
}

// PushFront 在队列前端插入元素
func (d *Deque) PushFront(value interface{}) {
	d.data.PushFront(value)
}

// PushBack 在队列后端插入元素
func (d *Deque) PushBack(value interface{}) {
	d.data.PushBack(value)
}

// PopFront 移除并返回队列前端的元素
func (d *Deque) PopFront() interface{} {
	front := d.data.Front()
	if front != nil {
		d.data.Remove(front)
		return front.Value
	}
	return nil
}

// PopBack 移除并返回队列后端的元素
func (d *Deque) PopBack() interface{} {
	back := d.data.Back()
	if back != nil {
		d.data.Remove(back)
		return back.Value
	}
	return nil
}

func main() {
	deque := NewDeque()

	// 测试 PushFront 和 PushBack
	deque.PushBack(1)
	deque.PushFront(2)
	deque.PushBack(3)
	deque.PushFront(4)

	// 测试 PopFront
	fmt.Println("Popped from front:", deque.PopFront()) // 输出: Popped from front: 4
	fmt.Println("Popped from front:", deque.PopFront()) // 输出: Popped from front: 2

	// 测试 PopBack
	fmt.Println("Popped from back:", deque.PopBack()) // 输出: Popped from back: 3
	fmt.Println("Popped from back:", deque.PopBack()) // 输出: Popped from back: 1

	// 测试空队列的情况
	fmt.Println("Popped from front on empty deque:", deque.PopFront()) // 输出: Popped from front on empty deque: <nil>
	fmt.Println("Popped from back on empty deque:", deque.PopBack())   // 输出: Popped from back on empty deque: <nil>

	// 再次测试 PushFront 和 PushBack
	deque.PushFront(5)
	deque.PushBack(6)

	// 测试 PeekFront 和 PeekBack
	fmt.Println("Front element:", deque.PopFront()) // 输出: Front element: 5
	fmt.Println("Back element:", deque.PopBack())   // 输出: Back element: 6
}
