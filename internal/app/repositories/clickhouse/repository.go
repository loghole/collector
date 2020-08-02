package clickhouse

import (
	"context"
	"database/sql"
	"time"

	"github.com/gadavy/tracing"
	"github.com/jmoiron/sqlx"

	"github.com/lissteron/loghole/collector/internal/app/domain"
)

const (
	insertLogsQuery = `INSERT INTO internal_logs_buffer 
            (time, 
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
             remote_ip) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
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

type EntryRepository struct {
	db     *sqlx.DB
	logger Logger

	period time.Duration
	queue  chan *domain.Entry
}

func NewEntryRepository(
	db *sqlx.DB,
	logger Logger,
	capacity int,
	period time.Duration,
) *EntryRepository {
	return &EntryRepository{
		db:     db,
		logger: logger,
		period: period,
		queue:  make(chan *domain.Entry, capacity),
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
	tx, stmt, err := r.getInsertEntryStmt()
	if err != nil {
		return err
	}

	var (
		entry  *domain.Entry
		count  uint
		active = true
		ticker = time.NewTicker(r.period)
	)

	defer ticker.Stop()

	for active {
		select {
		case <-ticker.C:
			if count == 0 {
				continue
			}

			if err := tx.Commit(); err != nil {
				r.logger.Error(ctx, err)
			}

			count = 0

			for {
				tx, stmt, err = r.getInsertEntryStmt()
				if err != nil {
					r.logger.Error(ctx, err)
					continue
				}

				break
			}
		case entry, active = <-r.queue:
			if !active {
				break
			}

			if err := r.insertEntry(stmt, entry); err != nil {
				r.logger.Error(ctx, err)
				continue
			}

			count++
		}
	}

	return tx.Commit()
}

func (r *EntryRepository) getInsertEntryStmt() (tx *sql.Tx, stmt *sql.Stmt, err error) {
	tx, err = r.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	stmt, err = tx.Prepare(insertLogsQuery)
	if err != nil {
		return nil, nil, err
	}

	return tx, stmt, nil
}

func (r *EntryRepository) insertEntry(stmt *sql.Stmt, entry *domain.Entry) (err error) {
	_, err = stmt.Exec(entry.Time, entry.Time, entry.Time.UnixNano(), entry.Namespace, entry.Source,
		entry.Host, entry.Level, entry.TraceID, entry.Message, string(entry.Params), entry.StringKey,
		entry.StringVal, entry.FloatKey, entry.FloatVal, entry.BuildCommit, entry.ConfigHash, entry.RemoteIP)

	return err
}
