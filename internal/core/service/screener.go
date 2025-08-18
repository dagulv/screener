package service

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
)

type Screener struct {
	store port.Screener
}

func NewScreener(store port.Screener) Screener {
	return Screener{
		store: store,
	}
}

func (s Screener) CountMagicRanks(ctx context.Context, filters domain.MagicRankFilter) (int, error) {
	return s.store.CountMagicRanks(ctx, filters)
}

func (s Screener) IterateMagicRanks(ctx context.Context, filters domain.MagicRankFilter) iter.Seq2[*domain.MagicRank, error] {
	return s.store.IterateMagicRanks(ctx, filters)
}

func (s Screener) CountScreener(ctx context.Context, filters domain.ScreenerFilter) (int, error) {
	return s.store.CountScreener(ctx, filters)
}

func (s Screener) IterateScreener(ctx context.Context, filters domain.ScreenerFilter) iter.Seq2[*domain.Screener, error] {
	if filters.OrderBy != "name" {
		var hasColumn bool
		for _, c := range filters.Columns {
			if c == filters.OrderBy {
				hasColumn = true
			}
		}

		if !hasColumn {
			filters.OrderBy = "name"
		}
	}

	return s.store.IterateScreener(ctx, filters)
}
