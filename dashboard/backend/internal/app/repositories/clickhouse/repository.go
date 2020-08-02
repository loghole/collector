package clickhouse

import (
	"context"

	"github.com/jmoiron/sqlx"
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

type Repository struct {
	db     *sqlx.DB
	logger Logger
}

func NewRepository(db *sqlx.DB, logger Logger) *Repository {
	return &Repository{db: db, logger: logger}
}
