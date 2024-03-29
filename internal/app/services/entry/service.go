package entry

import (
	"context"

	"github.com/lissteron/simplerr"
	"github.com/loghole/tracing"
	"github.com/loghole/tracing/tracelog"

	"github.com/loghole/collector/internal/app/codes"
	"github.com/loghole/collector/internal/app/domain"
)

type Storage interface {
	Ping(ctx context.Context) error
	StoreEntryList(ctx context.Context, list []*domain.Entry) (err error)
}

type Service struct {
	storage Storage
	logger  tracelog.Logger
}

func NewService(storage Storage, logger tracelog.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}

func (s *Service) Ping(ctx context.Context) error {
	defer tracing.ChildSpan(&ctx).Finish()

	if err := s.storage.Ping(ctx); err != nil {
		s.logger.Errorf(ctx, "ping db failed: %v", err)

		return simplerr.WrapWithCode(err, simplerr.InternalCode(codes.DatabaseError), "ping db failed")
	}

	return nil
}

func (s *Service) StoreItem(ctx context.Context, remoteIP string, data []byte) (err error) {
	defer tracing.ChildSpan(&ctx).Finish()

	entry, err := s.parseEntryItem(ctx, data)
	if err != nil {
		s.logger.Errorf(ctx, "parse entry item failed: %v", err)

		return simplerr.WrapWithCode(err, simplerr.InternalCode(codes.UnmarshalError), "parse json failed")
	}

	entry.SetRemoteIP(remoteIP)

	err = s.storage.StoreEntryList(ctx, []*domain.Entry{entry})
	if err != nil {
		s.logger.Errorf(ctx, "store entry list failed: %v", err)

		return simplerr.WrapWithCode(err, simplerr.InternalCode(codes.DatabaseError), "store failed")
	}

	return nil
}

func (s *Service) StoreList(ctx context.Context, remoteIP string, data []byte) (err error) {
	defer tracing.ChildSpan(&ctx).Finish()

	list, err := s.parseEntryList(ctx, data)
	if err != nil {
		s.logger.Errorf(ctx, "parse entry list failed: %v", err)

		return simplerr.WrapWithCode(err, simplerr.InternalCode(codes.UnmarshalError), "parse json failed")
	}

	list.SetRemoteIP(remoteIP)

	err = s.storage.StoreEntryList(ctx, list)
	if err != nil {
		s.logger.Errorf(ctx, "store entry list failed: %v", err)

		return simplerr.WrapWithCode(err, simplerr.InternalCode(codes.DatabaseError), "store failed")
	}

	return nil
}

func (s *Service) parseEntryItem(ctx context.Context, data []byte) (*domain.Entry, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	entry := &domain.Entry{}

	if err := entry.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *Service) parseEntryList(ctx context.Context, data []byte) (domain.EntryList, error) {
	defer tracing.ChildSpan(&ctx).Finish()

	list := domain.EntryList{}

	if err := list.UnmarshalJSON(data); err != nil {
		return nil, err
	}

	return list, nil
}
