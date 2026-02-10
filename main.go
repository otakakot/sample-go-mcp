package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-go-mcp",
		Title:   "Sample Go MCP Tool",
		Version: "v0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "say-hi", Description: "Standard Input/Output management tool"}, SayHi)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		panic(err)
	}
}

type Input struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type Output struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input Input,
) (
	*mcp.CallToolResult,
	Output,
	error,
) {
	return nil, Output{Greeting: "Hi " + input.Name}, nil
}
