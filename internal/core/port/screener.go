package port

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
)

type Screener interface {
	IterateScreener(ctx context.Context, filter domain.ScreenerFilter) iter.Seq2[*domain.Screener, error]
	CountScreener(ctx context.Context, filter domain.ScreenerFilter) (int, error)
	IterateMagicRanks(ctx context.Context, filter domain.MagicRankFilter) iter.Seq2[*domain.MagicRank, error]
	CountMagicRanks(ctx context.Context, filter domain.MagicRankFilter) (int, error)
}
