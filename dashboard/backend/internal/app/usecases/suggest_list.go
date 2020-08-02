package usecases

import (
	"context"
	"errors"

	"github.com/gadavy/tracing"
	"github.com/lissteron/simplerr"

	"github.com/lissteron/loghole/dashboard/internal/app/codes"
	"github.com/lissteron/loghole/dashboard/internal/app/domain"
	"github.com/lissteron/loghole/dashboard/internal/app/usecases/interfaces"
)

var ErrInvalidSuggestType = errors.New("invalid suggest type")

type ListSuggestIn struct {
	Type  string
	Value string
}

type ListSuggest struct {
	storage interfaces.Storage
	logger  interfaces.Logger
}

func NewListSuggest(
	storage interfaces.Storage,
	logger interfaces.Logger,
) *ListSuggest {
	return &ListSuggest{
		storage: storage,
		logger:  logger,
	}
}

func (l *ListSuggest) Do(ctx context.Context, input *ListSuggestIn) (result []string, err error) {
	defer tracing.ChildSpan(&ctx).Finish()

	switch input.Type {
	case domain.SuggestLevel:
		result, err = l.storage.ListLevelSuggest(ctx, input.Value)
	case domain.SuggestHost:
		result, err = l.storage.ListHostSuggest(ctx, input.Value)
	case domain.SuggestSource:
		result, err = l.storage.ListSourceSuggest(ctx, input.Value)
	case domain.SuggestNamespace:
		result, err = l.storage.ListNamespaceSuggest(ctx, input.Value)
	default:
		l.logger.Errorf(ctx, "invalid suggest type: %v", input.Type)
		return nil, simplerr.WrapWithCode(ErrInvalidSuggestType, codes.InvalidSuggestType, "invalid type")
	}

	if err != nil {
		l.logger.Errorf(ctx, "get suggest list failed: %v", err)
		return nil, simplerr.WrapWithCode(err, codes.DatabaseError, "get suggest list failed")
	}

	return result, nil
}
