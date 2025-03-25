package clientredis

import (
	"github.com/redis/go-redis/v9"
	"github.com/vtievsky/golibs/runtime/redisotel"
)

type ClientOpts struct {
	URL string
}

type Client struct {
	redis.UniversalClient
}

func New(opts *ClientOpts) (*Client, error) {
	r, err := redisotel.NewUniversalClient(
		&redis.UniversalOptions{ //nolint:exhaustruct
			Addrs:      []string{opts.URL},
			ClientName: "auth-id",
			DB:         0,
			PoolSize:   20, //nolint:mnd
		},
	)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return &Client{
		UniversalClient: r,
	}, nil
}
