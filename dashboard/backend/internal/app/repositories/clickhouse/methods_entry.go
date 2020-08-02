package clickhouse

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/gadavy/tracing"

	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/repositories/clickhouse/models"
)

func (r *Repository) ListEntry(ctx context.Context, input *domain.Query) ([]*domain.Entry, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	query, args, err := buildListEntryQuery(ctx, input)
	if err != nil {
		return nil, err
	}

	r.logger.Debugf(ctx, "query = %s", query)

	var dest []*models.Entry

	if err := r.db.SelectContext(ctx, &dest, query, args...); err != nil {
		return nil, err
	}

	result := make([]*domain.Entry, 0, len(dest))

	for _, val := range dest {
		result = append(result, val.ToDomain())
	}

	return result, nil
}

func buildListEntryQuery(
	ctx context.Context,
	input *domain.Query,
) (query string, args []interface{}, err error) {
	defer tracing.ChildSpan(&ctx).Finish()

	builder := squirrel.Select("time", "nsec", "namespace", "source", "host", "level",
		"trace_id", "message", "remote_ip", "params", "build_commit", "config_hash").From("internal_logs_buffer")

	for _, param := range input.Params {
		if param.IsTypeJSON() {
			builder = builder.Where(models.JSONParamFromDomain(param))
			continue
		}

		builder = builder.Where(models.ColumnParamFromDomain(param))
	}

	return builder.OrderBy("time DESC").
		Suffix(fmt.Sprintf("LIMIT %d, %d", input.Offset, input.Limit)).
		PlaceholderFormat(squirrel.Question).
		ToSql()
}
