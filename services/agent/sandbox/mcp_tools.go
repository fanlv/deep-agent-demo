package sandbox

import (
	"context"
	"fmt"
	"net/url"

	emcp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

func (c Client) GetMCPTools(ctx context.Context) ([]tool.BaseTool, error) {
	mcpURL, err := url.JoinPath(c.baseURL, "mcp")
	if err != nil {
		return nil, fmt.Errorf("failed to join URL path: %w", err)
	}

	fmt.Printf("[GetMCPTools] connecting to MCP server: %s\n", mcpURL)

	cli, err := client.NewStreamableHttpClient(mcpURL,
		transport.WithHTTPHeaders(map[string]string{
			"X-Session-ID": c.sessionID},
		// "X-Workspace": c.Ctx.Workspace,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP client: %w", err)
	}

	fmt.Printf("[GetMCPTools] starting MCP client...\n")
	err = cli.Start(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start MCP client (check if server is running at %s): %w", mcpURL, err)
	}

	fmt.Printf("[GetMCPTools] MCP client started, initializing...\n")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "deepagent-mcp-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MCP client: %w", err)
	}

	fmt.Printf("[GetMCPTools] MCP initialized, getting tools...\n")
	ts, err := emcp.GetTools(ctx, &emcp.Config{Cli: cli})
	if err != nil {
		return nil, fmt.Errorf("failed to get MCP tools: %w", err)
	}
	fmt.Printf("[GetMCPTools] got %d tools\n", len(ts))

	// fmt.Println("=== Available MCP Tools ===")
	// for i, mcpTool := range ts {
	// 	info, err := mcpTool.Info(ctx)
	// 	if err != nil {
	// 		log.Printf("Failed to get tool info: %v", err)
	// 		continue
	// 	}
	// 	fmt.Printf("%d. Name: %s\n", i+1, info.Name)
	// 	// fmt.Printf("%d. Description: %s\n", i+1, info.Desc)
	// }
	// fmt.Println("============================")

	return ts, nil
}
