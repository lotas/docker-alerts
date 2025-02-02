package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/system"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	return &Client{cli: cli}, nil
}

func (c *Client) Info(ctx context.Context) (system.Info, string, error) {
	info, err := c.cli.Info(ctx)
	if err != nil {
		return info, "", err
	}

	serverInfo := fmt.Sprintf("Docker version: %v\n", info.ServerVersion)
	serverInfo += fmt.Sprintf("Docker host: %v\n", info.Name)
	serverInfo += fmt.Sprintf("Type: %v\n", info.OSType)
	serverInfo += fmt.Sprintf("Architecture: %v\n", info.Architecture)
	serverInfo += fmt.Sprintf("CPUs: %v\n", info.NCPU)
	serverInfo += fmt.Sprintf("Memory: %v MB\n", info.MemTotal/1024/1024)

	return info, serverInfo, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
