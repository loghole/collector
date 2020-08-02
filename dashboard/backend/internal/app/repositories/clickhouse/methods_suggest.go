package clickhouse

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/gadavy/tracing"

	"github.com/lissteron/loghole/dashboard/internal/app/repositories/clickhouse/tools"
)

func (r *Repository) ListNamespaceSuggest(ctx context.Context, value string) ([]string, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	return r.listSuggest(ctx, "namespace", value)
}

func (r *Repository) ListSourceSuggest(ctx context.Context, value string) ([]string, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	return r.listSuggest(ctx, "source", value)
}

func (r *Repository) ListHostSuggest(ctx context.Context, value string) ([]string, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	return r.listSuggest(ctx, "host", value)
}

func (r *Repository) ListLevelSuggest(ctx context.Context, value string) ([]string, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	return r.listSuggest(ctx, "level", value)
}

func (r *Repository) listSuggest(ctx context.Context, column, value string) ([]string, error) {
	builder := squirrel.
		Select(column).
		From("internal_logs_buffer").
		GroupBy(column)

	if value != "" {
		builder = builder.Where(squirrel.Like{column: tools.CreateLike(value)})
	}

	query, args, err := builder.PlaceholderFormat(squirrel.Question).ToSql()
	if err != nil {
		return nil, err
	}

	var dest []string

	if err := r.db.SelectContext(ctx, &dest, query, args...); err != nil {
		return nil, err
	}

	return dest, nil
}
