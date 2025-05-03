package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"os"
)

func main() {
	//新建mcp 服务
	mcpServer := server.NewMCPServer("ziyi Mcp Server", "1.0.0")
	//新增一个mcp server对外暴露的静态资源（固定URI），比如：对外暴露项目的操作手册以及部署方式等，就可以通过静态资源README文件的方式来告诉大模型
	resource := mcp.NewResource(
		"docs://readme",
		"项目的README文件",
		mcp.WithResourceDescription("这是一个项目的README文件"),
		mcp.WithMIMEType("text/markdown"),
	)
	// 添加对静态资源的处理器
	mcpServer.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("README.md")
		if err != nil {
			return nil, err
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})
	//给该mcp server添加计算能力(Tools)
	// 1. 描述该工具，以及调用该工具调用的参数机器含义
	calculatorTool := mcp.NewTool("calculate",
		mcp.WithDescription("进行基础的数学运算"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)
	// 2. 实现工具的具体处理逻辑，类比大模型function_calling中的func部分
	mcpServer.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op := request.Params.Arguments["operation"].(string)
		x := request.Params.Arguments["x"].(float64)
		y := request.Params.Arguments["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return nil, errors.New("Division by zero is not allowed")
			}
			result = x / y
		}

		return mcp.FormatNumberResult(result), nil
	})

	// 给mcpServer添加提示词模版
	mcpServer.AddPrompt(mcp.NewPrompt("打招呼",
		mcp.WithPromptDescription("A friendly greeting prompt"),
		mcp.WithArgument("name",
			mcp.ArgumentDescription("Name of the person to greet"),
		),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := request.Params.Arguments["name"]
		if name == "" {
			name = "friend"
		}

		return mcp.NewGetPromptResult(
			"A friendly greeting",
			[]mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleAssistant,
					mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
				),
			},
		), nil
	})
}
