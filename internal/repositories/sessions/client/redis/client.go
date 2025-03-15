package clientredis

import "github.com/redis/go-redis/v9"

type ClientOpts struct {
	URL string
}

type Client struct {
	redis.UniversalClient
}

func New(opts *ClientOpts) *Client {
	return &Client{
		UniversalClient: redis.NewUniversalClient(
			&redis.UniversalOptions{ //nolint:exhaustruct
				Addrs:      []string{opts.URL},
				ClientName: "auth-id",
				DB:         0,
				PoolSize:   20, //nolint:mnd
			},
		),
	}
}
