package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/otakakot/sample-go-mcp/internal/tool"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "noa",
		Title:   "Notion Database Agent",
		Version: "v0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "stdio", Description: "Standard Input/Output management tool"}, tool.SayHi)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		panic(err)
	}
}
