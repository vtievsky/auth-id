package clienttarantool

import (
	"fmt"
	"net/url"
	"time"

	"github.com/tarantool/go-tarantool"
)

type Tuple []any

type ClientOpts struct {
	URL       string
	RateLimit uint
}

type Client struct {
	*tarantool.Connection
}

func New(opts *ClientOpts) (*Client, error) {
	databaseURL, err := url.Parse(opts.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL | %w", err)
	}

	c, err := tarantool.Connect(databaseURL.Host, tarantool.Opts{
		Auth:                 0,
		Dialer:               nil,
		Timeout:              time.Second,
		Reconnect:            time.Millisecond * 100, //nolint:mnd
		MaxReconnects:        3,                      //nolint:mnd
		User:                 "",
		Pass:                 "",
		RateLimit:            0,
		RLimitAction:         0,
		Concurrency:          0,
		SkipSchema:           false,
		Notify:               make(chan<- tarantool.ConnEvent),
		Handle:               nil,
		Logger:               nil,
		Transport:            "",
		Ssl:                  tarantool.SslOpts{},      //nolint:exhaustruct
		RequiredProtocolInfo: tarantool.ProtocolInfo{}, //nolint:exhaustruct
	})
	if err != nil {
		return nil, fmt.Errorf("failed create tarantool connection | %w", err)
	}

	return &Client{
		Connection: c,
	}, nil
}
