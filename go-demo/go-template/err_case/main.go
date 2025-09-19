package main

import (
	"fmt"
	"os"
	"text/template"
)

// 错误结构
type BusinessError struct {
	Code    int
	Message string
}

func (e BusinessError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// 在模板中抛出错误的核心函数
func throwError(code int, message string) (interface{}, error) {
	return nil, BusinessError{
		Code:    code,
		Message: message,
	}
}

func main() {
	tmpl := `
{{/* 条件性抛出错误 */}}
{{$age := 15}}
{{ if gt $age 18}}
	已成年
{{else}}
	{{throwError 1001 "年龄不满足"}}
{{end}}
`

	funcMap := template.FuncMap{
		"throwError": throwError,
	}

	parsedTmpl, err := template.New("simple_error").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		fmt.Printf("模板解析错误: %v\n", err)
		return
	}

	err = parsedTmpl.Execute(os.Stdout, nil)
	if err != nil {
		fmt.Printf("执行模板错误: %v\n", err)
	}
}
