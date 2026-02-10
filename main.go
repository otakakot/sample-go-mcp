package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-go-mcp",
		Title:   "Sample Go MCP Tool",
		Version: "v0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "say-hi",
		Description: "Greets the user with a friendly message",
	}, SayHi)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "elicitation",
		Description: "Echoes back the input message",
	}, Elicitation)

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		panic(err)
	}
}

type SayHiInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type SayHiOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func SayHi(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SayHiInput,
) (
	*mcp.CallToolResult,
	SayHiOutput,
	error,
) {
	return nil, SayHiOutput{Greeting: "Hi " + input.Name}, nil
}

type ElicitationInput struct{}

type ElicitationOutput struct {
	ElicitedMessage string `json:"elicited_message" jsonschema:"the elicited message"`
}

func Elicitation(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ElicitationInput,
) (
	*mcp.CallToolResult,
	ElicitationOutput,
	error,
) {
	elicit, err := req.Session.Elicit(ctx, &mcp.ElicitParams{
		Message: "Please enter your message:",
		RequestedSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"message": map[string]string{
					"type":        "string",
					"description": "the message to tell to the AI agent",
				},
			},
			"required": []string{"message"},
		},
	})
	if err != nil {
		return nil, ElicitationOutput{}, fmt.Errorf("elicitation failed: %w", err)
	}

	message, ok := elicit.Content["message"].(string)
	if !ok {
		return nil, ElicitationOutput{}, fmt.Errorf("invalid elicitation result")
	}

	return nil, ElicitationOutput{ElicitedMessage: message}, nil
}
