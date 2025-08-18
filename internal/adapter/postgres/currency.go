package postgres

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
	"github.com/webmafia/pg"
)

type currencyStore struct {
	db
}

func NewCurrency(pool *pg.DB) port.Currency {
	return currencyStore{
		db: db{pool},
	}
}

// CreateCurrency implements port.Currency
func (s currencyStore) CreateCurrency(ctx context.Context, currency *domain.IDAndName) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("id", currency.ID).
		Value("name", currency.Name)

	_, err = s.db.InsertValues(ctx, Currency, vals)

	return
}

// ReadCurrency implements port.Currency
func (s currencyStore) ReadCurrency(ctx context.Context, currency *domain.IDAndName) (err error) {
	c := Currency.Alias("c")

	row := s.db.QueryRow(ctx, `
		SELECT
			c.id,
			c.name
		FROM %T
		WHERE %c
	`, c, Currency, Currency, pg.Eq(c.Col("id"), currency.ID))

	err = row.Scan(
		&currency.ID,
		&currency.Name,
	)

	return
}

// UpdateCurrency implements port.Currency
func (s currencyStore) UpdateCurrency(ctx context.Context, currency *domain.IDAndName) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("name", currency.Name)

	_, err = s.db.UpdateValues(ctx, Currency, vals, pg.Eq("id", currency.ID))

	return
}

// DeleteCurrency implements port.Currency
func (s currencyStore) DeleteCurrency(ctx context.Context, currencyId xid.ID) (err error) {
	_, err = s.db.Delete(ctx, Currency, pg.Eq("id", currencyId))
	return
}

// IterateCurrencys implements port.Currency
func (s currencyStore) IterateCurrencies(ctx context.Context, filters domain.IDAndNameFilter) iter.Seq2[*domain.IDAndName, error] {
	return func(yield func(*domain.IDAndName, error) bool) {
		c := Currency.Alias("c")
		cond := pg.And()

		rows, err := s.db.Query(ctx, `
			SELECT
				c.id,
				c.name
			FROM %T
			WHERE %c
		`, c, cond)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var currency domain.IDAndName

			if err = rows.Scan(
				&currency.ID,
				&currency.Name,
			); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&currency, nil) {
				return
			}
		}
	}
}
