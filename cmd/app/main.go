package main

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "app",
		Version: "0.0.1",
	}, nil)

	server.AddResource(&mcp.Resource{
		Meta:        mcp.Meta{},
		Annotations: &mcp.Annotations{},
		Description: "A greeting resource",
		MIMEType:    "text/html;profile=mcp-app",
		Name:        "greeting",
		Size:        0,
		Title:       "Greeting",
		URI:         "ui://app/greeting",
		Icons:       []mcp.Icon{},
	}, Resource)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "hello_world",
		Description: "Display a Hello World greeting with optional interactive UI",
		Meta: mcp.Meta{
			"ui": map[string]any{
				"resourceUri": "ui://app/greeting",
				"visibility":  []string{"model", "app"},
			},
		},
		Annotations: &mcp.ToolAnnotations{},
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"name": map[string]any{
					"type":        "string",
					"description": "the name of the person to greet",
				},
			},
			"required": []string{"name"},
		},
		OutputSchema: nil,
		Title:        "hello_world",
		Icons:        []mcp.Icon{},
	}, SayHi)

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

func Resource(
	ctx context.Context,
	req *mcp.ReadResourceRequest,
) (
	*mcp.ReadResourceResult,
	error,
) {
	return &mcp.ReadResourceResult{
		Meta: mcp.Meta{},
		Contents: []*mcp.ResourceContents{
			{
				URI:      "ui://app/greeting",
				MIMEType: "text/html;profile=mcp-app",
				Text:     ui,
				Meta: mcp.Meta{
					"ui":            "",
					"prefersBorder": false,
				},
			},
		},
	}, nil
}

const ui = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Hello World</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      max-width: 400px;
      margin: 20px auto;
      padding: 20px;
    }
    .greeting {
      font-size: 24px;
      margin-bottom: 20px;
    }
    input {
      width: 100%;
      padding: 8px;
      margin-bottom: 10px;
      border: 1px solid #ccc;
      border-radius: 4px;
    }
    button {
      background: #007acc;
      color: white;
      padding: 10px 20px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
    button:hover {
      background: #005aa3;
    }
    .status {
      margin-top: 10px;
      color: #666;
    }
  </style>
</head>
<body>
  <div class="greeting" id="greeting">Hello, World!</div>
  
  <input type="text" id="nameInput" placeholder="Enter your name" value="World">
  <button onclick="greet()">Greet</button>
  
  <div class="status" id="status"></div>

  <script>
    let nextRequestId = 1;
    const pendingRequests = new Map();

    function sendRequest(method, params) {
      const id = nextRequestId++;
      const request = {
        jsonrpc: '2.0',
        id: id,
        method: method,
        params: params || {}
      };
      
      return new Promise((resolve, reject) => {
        pendingRequests.set(id, { resolve, reject });
        window.parent.postMessage(request, '*');
      });
    }

    window.addEventListener('message', (event) => {
      const message = event.data;
      if (!message || message.jsonrpc !== '2.0') return;
      
      if (message.id !== undefined && (message.result !== undefined || message.error)) {
        const pending = pendingRequests.get(message.id);
        if (pending) {
          pendingRequests.delete(message.id);
          if (message.error) {
            pending.reject(new Error(message.error.message || 'Unknown error'));
          } else {
            pending.resolve(message.result);
          }
        }
      }
    });

    async function greet() {
      const name = document.getElementById('nameInput').value.trim() || 'World';
      document.getElementById('status').textContent = 'Calling tool...';
      
      try {
        const result = await sendRequest('tools/call', {
          name: 'hello_world',
          arguments: { name: name }
        });
        
        if (result?.content && result.content.length > 0) {
          document.getElementById('greeting').textContent = result.content[0].text;
          document.getElementById('status').textContent = 'Success!';
        } else {
          throw new Error('No response');
        }
      } catch (error) {
        document.getElementById('status').textContent = 'Error: ' + error.message;
      }
    }
  </script>
</body>
</html>`
