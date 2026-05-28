package Strategy

import (
	"context"

	"github.com/docker/docker/client"
)

// struct is used to define the common fields for all database strategies
// this is to ensure that all database strategies have the same fields and can be used interchangeably
// this is the same define the required fields for all concrete implementations of the DatabaseStrategyInterface
type DatabaseStrategy struct {
	InstanceName      string
	HostPort          string
	LocalDatabaseName string
	LocalDatabaseDir  string
	Username          string
	Password          string
}

// interface is used to define the common methods for all database strategies
type DatabaseStrategyInterface interface {
	Deploy(ctx context.Context, cli *client.Client, cfg DatabaseStrategy) (string, error)
	Connect(ctx context.Context, cli *client.Client, cfg DatabaseStrategy) error
	Close(ctx context.Context) error
}
