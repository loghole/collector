package clickhouse

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/loghole/database"
	"github.com/loghole/gorand"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracelog"

	"github.com/loghole/collector/internal/app/domain"
)

const (
	insertLogsQuery = `INSERT INTO internal_logs_buffer (
		time, 
		date, 
		nsec, 
		namespace, 
		source, 
		host, 
		level, 
		trace_id,
		message, 
		params, 
		params_string.keys, 
		params_string.values, 
		params_float.keys, 
		params_float.values, 
		build_commit, 
		config_hash,
		remote_ip,
		row_id) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
)

type EntryRepository struct {
	db     *database.DB
	logger tracelog.Logger

	period time.Duration
	queue  chan *domain.Entry

	rand rand.Source64
}

func NewEntryRepository(
	db *database.DB,
	logger tracelog.Logger,
	capacity int,
	period time.Duration,
) *EntryRepository {
	return &EntryRepository{
		db:     db,
		logger: logger,
		period: period,
		queue:  make(chan *domain.Entry, capacity),
		rand:   rand.Source64(rand.New(gorand.NewSource(time.Now().UnixNano()))), //nolint:gosec // need pseudo random
	}
}

func (r *EntryRepository) Ping(ctx context.Context) error {
	defer tracing.ChildSpan(&ctx).Finish()

	return r.db.PingContext(ctx)
}

func (r *EntryRepository) Run(ctx context.Context) error {
	return r.storeEntryChan(ctx)
}

func (r *EntryRepository) Stop() {
	close(r.queue)
}

func (r *EntryRepository) StoreEntryList(ctx context.Context, list []*domain.Entry) (err error) {
	defer tracing.ChildSpan(&ctx).Finish()

	for _, entry := range list {
		r.queue <- entry
	}

	return nil
}

func (r *EntryRepository) storeEntryChan(ctx context.Context) error {
	var (
		entry  *domain.Entry
		active = true
		ticker = time.NewTicker(r.period)

		cache = make([]*domain.Entry, 0)
	)

	defer ticker.Stop()

	for active {
		select {
		case <-ticker.C:
			if len(cache) == 0 {
				continue
			}

			if err := r.insertEntryList(ctx, cache); err != nil {
				r.logger.Errorf(ctx, "insert entry list: %v", err)
			}

			cache = make([]*domain.Entry, 0, len(cache))
		case entry, active = <-r.queue:
			if !active {
				break
			}

			cache = append(cache, entry)
		}
	}

	if len(cache) > 0 {
		if err := r.insertEntryList(ctx, cache); err != nil {
			r.logger.Errorf(ctx, "insert entry list: %v", err)
		}
	}

	return nil
}

func (r *EntryRepository) insertEntryList(ctx context.Context, cache []*domain.Entry) error {
	err := r.db.RunTxx(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		stmt, err := tx.Prepare(insertLogsQuery)
		if err != nil {
			return fmt.Errorf("prepare stmt: %w", err)
		}

		defer func() {
			if err := stmt.Close(); err != nil {
				r.logger.Errorf(ctx, "stmt close: %v", err)
			}
		}()

		for _, entry := range cache {
			if _, err := stmt.Exec(
				entry.Time,
				entry.Time,
				entry.Time.UnixNano(),
				entry.Namespace,
				entry.Source,
				entry.Host,
				entry.Level,
				entry.TraceID,
				entry.Message,
				string(entry.Params),
				entry.StringKey,
				entry.StringVal,
				entry.FloatKey,
				entry.FloatVal,
				entry.BuildCommit,
				entry.ConfigHash,
				entry.RemoteIP,
				r.rand.Uint64(),
			); err != nil {
				return fmt.Errorf("insert: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("transaction: %w", err)
	}

	return nil
}
