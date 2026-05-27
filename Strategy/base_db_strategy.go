package Strategy

import "context"


type DatabaseStrategy struct {
	InstanceName string;
	HostPort string;
	LocalDatabaseName string;
}


type DatabaseStrategyInterface interface {
	Connect(ctx context.Context) error;
	Close(ctx context.Context) error;
	ExecuteQuery(ctx context.Context, query string) (interface{}, error);
}


