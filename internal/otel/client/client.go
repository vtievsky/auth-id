package otelclient

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OtelOpts struct {
	URL string
}

type OtelClient struct {
	*grpc.ClientConn
}

func New(opts *OtelOpts) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		opts.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, nil
}
