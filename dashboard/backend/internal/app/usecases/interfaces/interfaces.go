package interfaces

import (
	"context"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, template string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, template string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, template string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, template string, args ...interface{})
}

type Storage interface {
	// Entry methods
	ListEntry(ctx context.Context, query *domain.Query) ([]*domain.Entry, error)

	// Suggestion methods
	ListNamespaceSuggest(ctx context.Context, value string) ([]string, error)
	ListSourceSuggest(ctx context.Context, value string) ([]string, error)
	ListHostSuggest(ctx context.Context, value string) ([]string, error)
	ListLevelSuggest(ctx context.Context, value string) ([]string, error)
}
