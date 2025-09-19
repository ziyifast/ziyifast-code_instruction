package main

import (
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	//案例：基础用法
	//case_basic_use()

	//案例：注释使用
	//case_comment_use()

	//案例：获取结构体属性
	//case_get_field()

	//案例：条件判断
	//case_condition()

	//案例：比较判断
	case_comparison()

	//案例：循环
	//case_loop()

	//案例：定义变量
	//case_define_variable()

	//案例：管道
	//case_pipe()

	//案例：自定义函数
	//case_custom_function()

	//案例：模版嵌套
	//case_template_nesting()

	//案例：空行去除
	//case_remove_empty_line()

}

// 移除空行
// {{-：去除左侧的空白符（包括空格、制表符和换行符）
// -}}：去除右侧的空白符
// {{- -}}：去除两侧的空白符
func case_remove_empty_line() {
	//不去除空行的模版
	tmpl := `
{{range .Lines}}
{{.}}
{{end}}`
	//去除空行的模版
	tmpl = `
{{- range .Lines -}}
{{- . -}}
{{- end -}}`
	data := map[string]interface{}{
		"Lines": []string{"", "line1", "", "line2", ""},
	}
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 模版嵌套
func case_template_nesting() {
	tmpl := `
{{define "header"}}
================================
{{.Title}}
================================
{{end}}

{{define "content"}}
用户信息:
姓名: {{.User.Name}}
邮箱: {{.User.Email}}
{{end}}

{{define "footer"}}
--------------------------------
报告生成时间: {{.Time}}
--------------------------------
{{end}}

{{template "header" .}}
{{template "content" .}}
{{template "footer" .}}`

	data := map[string]interface{}{
		"Title": "用户报告",
		"User": map[string]string{
			"Name":  "张三",
			"Email": "zhangsan@example.com",
		},
		"Time": "2023-01-01 12:00:00",
	}

	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 自定义函数
func case_custom_function() {
	funcMap := template.FuncMap{
		"formatTime": formatTime,
		"multiply":   multiply,
		"upper":      strings.ToUpper,
	}

	data := map[string]interface{}{
		"Name":     "张三",
		"Time":     time.Now(),
		"Price":    100,
		"Quantity": 3,
	}

	tmpl := `用户: {{.Name | upper}}
当前时间: {{.Time | formatTime}}
商品总价: {{multiply .Price .Quantity}} 元`

	t := template.Must(template.New("example").Funcs(funcMap).Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func multiply(a, b int) int {
	return a * b
}

// 管道
func case_pipe() {
	data := map[string]interface{}{
		"Name":  "zhangsan",
		"Text":  "hello world",
		"Items": []string{"item1", "item2", "item3"},
	}

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
	}

	tmpl := `原始名称: {{.Name}}
大写名称: {{.Name | upper}}
项目数量: {{len .Items}}
项目连接: {{join "," .Items}}

格式化: {{printf "Hello %s, you have %d items" .Name (len .Items)}}`

	t := template.Must(template.New("example").Funcs(funcMap).Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 定义变量
func case_define_variable() {
	data := map[string]interface{}{
		"Items": []string{"item1", "item2", "item3"},
		"Name":  "张三",
	}

	tmpl := `{{$name := .Name}}
{{$count := len .Items}}
用户: {{$name}}
项目数量: {{$count}}

{{range $index, $item := .Items}}
项目 {{$index}}: {{$item}}
{{end}}`

	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 循环
func case_loop() {
	data := map[string][]string{
		"Fruits": {"苹果", "香蕉", "橙子"},
	}

	tmpl := `水果列表:
{{range .Fruits}}
- {{.}}
{{end}}

带索引的列表:
{{range $index, $fruit := .Fruits}}
{{add $index 1}}. {{$fruit}}
{{end}}

空列表处理:
{{range .Vegetables}}
- {{.}}
{{else}}
没有蔬菜
{{end}}`

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	t := template.Must(template.New("example").Funcs(funcMap).Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func case_comparison() {
	user := map[string]interface{}{
		"name": "curry",
		"age":  20,
	}
	//比较判断
	//- eq: 等于 (==)
	//- ne: 不等于 (!=)
	//- lt: 小于 (<)
	//- le: 小于等于 (<=)
	//- gt: 大于 (>)
	//- ge: 大于等于 (>=)
	tmpl := `{{if ge .age 18}} 
				已成年
			 {{else}}
				未成年
			 {{end}}`
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}

// 条件判断
func case_condition() {
	data := map[string]interface{}{
		"Age":    25,
		"Name":   "张三",
		"Active": true,
	}
	// if、else条件判断
	tmpl := `{{if .Active}}
用户 {{.Name}} 处于活跃状态
{{else}}
用户 {{.Name}} 处于非活跃状态
{{end}}
{{if gt .Age 18}}
用户已成年
{{else}}
用户未成年
{{end}}`

	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 获取结构体属性
func case_get_field() {
	type Data struct {
		Name    string
		Hobbies []string
		Info    map[string]string
	}

	data := Data{
		Name:    "张三",
		Hobbies: []string{"读书", "游泳", "跑步"},
		Info:    map[string]string{"city": "北京", "job": "工程师"},
	}

	tmpl := `姓名: {{.Name}}
第一个爱好: {{index .Hobbies 0}}
城市: {{.Info.city}}`

	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

// 包含注释
func case_comment_use() {
	tmpl := `{{/* 这是一个注释 */}}
Hello, {{.Name}}!
{{/* 另一个注释 */}}`
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(os.Stdout, map[string]string{"Name": "世界"})
	if err != nil {
		panic(err)
	}
}

// 基础使用
func case_basic_use() {
	type Person struct {
		Name string
		Age  int
	}

	tmpl := `姓名: {{.Name}}, 年龄: {{.Age}}`
	t := template.Must(template.New("example").Parse(tmpl))

	person := Person{Name: "张三", Age: 30}
	err := t.Execute(os.Stdout, person)
	if err != nil {
		panic(err)
	}
}
