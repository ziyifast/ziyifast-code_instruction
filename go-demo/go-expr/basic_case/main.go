package main

import (
	"fmt"
	"github.com/expr-lang/expr"
)

func main() {
	BasicAddExpr()       //基础的加法表达式（不含变量）
	ContainsVarEnvCase() //包含变量的表达式
}

// ContainsVarEnvCase 包含变量的表达式
func ContainsVarEnvCase() {
	// 创建一个包含变量的环境
	env := map[string]interface{}{
		"apple":  5,
		"banana": 10,
	}

	// 编译一个使用环境中变量的表达式
	program, err := expr.Compile(`apple + banana`, expr.Env(env))
	if err != nil {
		panic(err)
	}

	// 运行表达式
	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	// 打印结果
	fmt.Println(output) // 输出 15
}

// BasicAddExpr 基础的加法表达式（不含变量）
func BasicAddExpr() {
	// 编译一个基础的加法表达式
	program, err := expr.Compile(`2 + 2`)
	if err != nil {
		panic(err)
	}

	// 运行编译后的表达式，并没有传入环境，因为这里不需要使用任何变量
	output, err := expr.Run(program, nil)
	if err != nil {
		panic(err)
	}

	// 打印结果
	fmt.Println(output) // 输出 4
}
