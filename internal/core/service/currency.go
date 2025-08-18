package service

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
)

type Currency struct {
	store port.Currency
}

func NewCurrency(store port.Currency) Currency {
	return Currency{
		store: store,
	}
}

func (s Currency) Create(ctx context.Context, currency *domain.IDAndName) (err error) {
	currency.ID = xid.New()

	return s.store.CreateCurrency(ctx, currency)
}

func (s Currency) Read(ctx context.Context, currency *domain.IDAndName) (err error) {
	return s.store.ReadCurrency(ctx, currency)
}

func (s Currency) Update(ctx context.Context, currency *domain.IDAndName) (err error) {
	return s.store.UpdateCurrency(ctx, currency)
}

func (s Currency) Delete(ctx context.Context, currencyId xid.ID) (err error) {
	return s.store.DeleteCurrency(ctx, currencyId)
}

func (s Currency) Iterate(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error] {
	return s.store.IterateCurrencies(ctx, filters)
}
