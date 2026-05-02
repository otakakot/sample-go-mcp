package main

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/otakakot/sample-go-mcp/internal/tool"
)

func main() {
	port := cmp.Or(os.Getenv("PORT"), "8080")

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "sample-go-mcp-sse",
		Version: "0.0.1",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{Name: "sse", Description: "Server-Sent Events management tool"}, tool.SayHi)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	http.Handle("/", handler)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	go func() {
		slog.Info("start server listen")

		if err := http.ListenAndServe(":"+port, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()
}
