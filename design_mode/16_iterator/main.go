package main

import "fmt"

/*
迭代器模式（Iterator Pattern）：遍历一个聚合对象的元素而无需暴露其内部实现。
迭代器模式提供了一种方式来顺序访问一个聚合对象中的各个元素，而不暴露聚合对象的内部实现细节。

实现细节：
1. 通过定义一个 Iterator 接口来实现迭代器模式。
2. 该Iterator接口需要定义两个方法：Next 和 HasNext。
  - Next 方法用于获取下一个元素
  - HasNext 方法用于判断是否还有下一个元素。

3. 将定义一个 Aggregate 接口和一个具体的聚合对象类型来实现迭代器模式。
*/
// Iterator 迭代器接口
type Iterator interface {
	Next() interface{}
	HasNext() bool
}

// 具体的聚合对象类型
type Numbers struct {
	numbers []int
}

func (n *Numbers) Iterator() Iterator {
	return &NumberIterator{
		numbers: n.numbers,
		index:   0,
	}
}

// NumberIterator 数字迭代器
type NumberIterator struct {
	numbers []int
	index   int
}

func (ni *NumberIterator) Next() interface{} {
	number := ni.numbers[ni.index]
	ni.index++
	return number
}

func (ni *NumberIterator) HasNext() bool {
	if ni.index >= len(ni.numbers) {
		return false
	}
	return true
}

func main() {
	numbers := &Numbers{
		numbers: []int{1, 2, 3, 4, 5},
	}
	iterator := numbers.Iterator()

	for iterator.HasNext() {
		fmt.Println(iterator.Next())
	}
}
