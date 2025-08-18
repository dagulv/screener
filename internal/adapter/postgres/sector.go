package postgres

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
	"github.com/webmafia/pg"
)

type sectorStore struct {
	db
}

func NewSector(pool *pg.DB) port.Sector {
	return sectorStore{
		db: db{pool},
	}
}

// CreateSector implements port.Sector
func (s sectorStore) CreateSector(ctx context.Context, sector *domain.IDAndName) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("id", sector.ID).
		Value("name", sector.Name)

	_, err = s.db.InsertValues(ctx, Sector, vals)

	return
}

// ReadSector implements port.Sector
func (s sectorStore) ReadSector(ctx context.Context, sector *domain.IDAndName) (err error) {
	se := Sector.Alias("s")

	row := s.db.QueryRow(ctx, `
		SELECT
			s.id,
			s.name
		FROM %T
		WHERE %c
	`, se, pg.Eq(se.Col("id"), sector.ID))

	err = row.Scan(
		&sector.ID,
		&sector.Name,
	)

	return
}

// UpdateSector implements port.Sector
func (s sectorStore) UpdateSector(ctx context.Context, sector *domain.IDAndName) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("name", sector.Name)

	_, err = s.db.UpdateValues(ctx, Sector, vals, pg.Eq("id", sector.ID))

	return
}

// DeleteSector implements port.Sector
func (s sectorStore) DeleteSector(ctx context.Context, sectorId xid.ID) (err error) {
	_, err = s.db.Delete(ctx, Sector, pg.Eq("id", sectorId))
	return
}

// IterateSectors implements port.Sector
func (s sectorStore) IterateSectors(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error] {
	return func(yield func(*domain.IDAndName, error) bool) {
		se := Sector.Alias("s")
		cond := pg.And()

		rows, err := s.db.Query(ctx, `
			SELECT
				s.id,
				s.name
			FROM %T
			WHERE %c
		`, se, cond)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var sector domain.IDAndName

			if err = rows.Scan(
				&sector.ID,
				&sector.Name,
			); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&sector, nil) {
				return
			}
		}
	}
}
