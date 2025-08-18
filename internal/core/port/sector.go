package port

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/rs/xid"
)

type Sector interface {
	Context

	CreateSector(ctx context.Context, sector *domain.IDAndName) error
	ReadSector(ctx context.Context, sector *domain.IDAndName) error
	UpdateSector(ctx context.Context, sector *domain.IDAndName) error
	DeleteSector(ctx context.Context, sectorId xid.ID) error
	IterateSectors(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error]
}
