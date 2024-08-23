package main

import (
	"github.com/agiledragon/gomonkey"
	_ "github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

// 案例一：一个Convey，一个用例
func TestEqualWithSingleTestCase(t *testing.T) {
	// test name：用例名称
	// t：需要传入*testing.T
	// func(){} 测试函数
	Convey("test name", t, func() {
		//1+1：断言
		//ShouldEqual：convey内置的断言
		//2：期望结果
		So(1+1, ShouldEqual, 2)
	})
}

// 案例二：多个Convey，多个用例（平铺写法）
func TestEqualWithMultipleTestCase(t *testing.T) {
	Convey("test add case", t, func() {
		So(1+1, ShouldEqual, 2)
	})
	Convey("test sub case", t, func() {
		So(1-1, ShouldEqual, 0)
	})
	Convey("test multi case", t, func() {
		So(1*1, ShouldNotEqual, -1)
	})
}

// 案例二：多个Convey，多个用例（嵌套写法）
func TestEqualWithMultipleTestCaseAndNested(t *testing.T) {
	Convey("test case", t, func() {
		Convey("test add case", func() {
			So(1+1, ShouldEqual, 2)
		})
		Convey("test sub case", func() {
			So(1-1, ShouldEqual, 0)
		})
		Convey("test multi case", func() {
			So(1*1, ShouldNotEqual, -1)
		})
	})
}

// 案例三：函数式断言
func TestFunctionalAssertion(t *testing.T) {
	Convey("test case", t, func() {
		So(add(1, 1), ShouldEqual, 2)
	})
}

func add(a, b int) int {
	return a + b
}

// 案例四：忽略Convey断言
// 忽略所有断言
func TestCaseSkipConvey(t *testing.T) {
	SkipConvey("test case", t, func() {
		So(add(1, 1), ShouldEqual, 2)
	})
}

// 忽略某些断言（SkipSo的断言将被忽略）
func TestCaseSkipSo(t *testing.T) {
	Convey("test case", t, func() {
		SkipSo(add(1, 1), ShouldEqual, 2)
		So(1-1, ShouldEqual, 0)
	})
}

// 案例五：定制断言
func TestCustomAssertion(t *testing.T) {
	Convey("test custom assert", t, func() {
		So(1+1, CustomAssertionWithRaiseMoney, 2)
	})
}

func CustomAssertionWithRaiseMoney(actual any, expected ...any) string {
	if actual == expected[0] {
		return ""
	} else {
		return "doesn't raise money"
	}
}

// 拓展：配合monkey打桩
var num = 10 //全局变量

func TestApplyGlobalVar(t *testing.T) {
	Convey("TestApplyGlobalVar", t, func() {
		Convey("change", func() {
			//模拟函数行为，给全局变量复制，在函数结束后直接通过reset恢复全局变量值
			patches := gomonkey.ApplyGlobalVar(&num, 150)
			defer patches.Reset()
			So(num, ShouldEqual, 150)
		})

		Convey("recover", func() {
			So(num, ShouldEqual, 10)
		})
	})
}

// 对函数进行打桩
func TestFunc(t *testing.T) {
	// mock 了 networkCompute()，返回了计算结果2
	patches := gomonkey.ApplyFunc(networkCompute, func(a, b int) (int, error) {
		return 2, nil
	})
	defer patches.Reset()
	sum, err := Compute(1, 2)
	println("expected %v, got %v", 2, sum)
	if sum != 2 || err != nil {
		t.Errorf("expected %v, got %v", 2, sum)
	}
}

func networkCompute(a, b int) (int, error) {
	// do something in remote computer
	c := a + b
	return c, nil
}

func Compute(a, b int) (int, error) {
	sum, err := networkCompute(a, b)
	return sum, err
}
