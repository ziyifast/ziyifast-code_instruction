package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"io"
	"net"
	"net/http"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"ip-mcp",
		"1.0.0",
	)

	// Add tool
	tool := mcp.NewTool("ip_query",
		mcp.WithDescription("query geo location of an IP address"),
		mcp.WithString("ip",
			mcp.Required(),
			mcp.Description("IP address to query"),
		),
	)

	// Add tool handler
	s.AddTool(tool, ipQueryHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func ipQueryHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ip, ok := request.Params.Arguments["ip"].(string)
	if !ok {
		return nil, errors.New("ip must be a string")
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return nil, errors.New("invalid IP address")
	}

	resp, err := http.Get("https://ip.rpcx.io/api/ip?ip=" + ip)
	if err != nil {
		return nil, fmt.Errorf("Error fetching IP information: %v", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}
	fmt.Printf("call ip: %s  Response body: %s\n", ip, string(data))
	return mcp.NewToolResultText(string(data)), nil
}
