package port

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/rs/xid"
)

type Currency interface {
	Context

	CreateCurrency(ctx context.Context, currency *domain.IDAndName) error
	ReadCurrency(ctx context.Context, currency *domain.IDAndName) error
	UpdateCurrency(ctx context.Context, currency *domain.IDAndName) error
	DeleteCurrency(ctx context.Context, currencyId xid.ID) error
	IterateCurrencies(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error]
}
