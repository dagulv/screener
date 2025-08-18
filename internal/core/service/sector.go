package service

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
)

type Sector struct {
	store port.Sector
}

func NewSector(store port.Sector) Sector {
	return Sector{
		store: store,
	}
}

func (s Sector) Create(ctx context.Context, sector *domain.IDAndName) (err error) {
	sector.ID = xid.New()

	return s.store.CreateSector(ctx, sector)
}

func (s Sector) Read(ctx context.Context, sector *domain.IDAndName) (err error) {
	return s.store.ReadSector(ctx, sector)
}

func (s Sector) Update(ctx context.Context, sector *domain.IDAndName) (err error) {
	return s.store.UpdateSector(ctx, sector)
}

func (s Sector) Delete(ctx context.Context, sectorId xid.ID) (err error) {
	return s.store.DeleteSector(ctx, sectorId)
}

func (s Sector) Iterate(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error] {
	return s.store.IterateSectors(ctx, filters)
}
