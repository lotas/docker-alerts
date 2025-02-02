package docker

import (
	"context"

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

func (c *Client) Info(ctx context.Context) (system.Info, error) {
  return c.cli.Info(ctx)
}

func (c *Client) Close() error {
	return c.cli.Close()
}
