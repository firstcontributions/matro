package store

const storeTmpl = `
package grpc

import (
	pool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

type {{ title .Name -}}Store struct {
	pool *pool.Pool
}

// New {{- title .Name -}}Store makes and keeps connection pool to given grpc server
// and return an instance of the client
func New {{- title .Name -}}Store(ctx context.Context, url string, initConnections, connectionCapacity, ttl int) (* {{ title .Name -}}Store, error) {
	pool, err := pool.New(
		func() (*grpc.ClientConn, error) {
			return grpc.Dial(url, grpc.WithInsecure())
		},
		initConnections,
		connectionCapacity,
		time.Duration(ttl)*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	return &{{ title .Name -}}Store{
		pool: pool,
	}, nil
}  

`
