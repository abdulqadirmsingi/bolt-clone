package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func Connect(ctx context.Context, addr, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &Client{Client: rdb}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
