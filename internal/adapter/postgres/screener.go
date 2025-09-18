package postgres

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/webmafia/pg"
)

type screenerStore struct {
	db
}

func NewScreener(pool *pg.DB) port.Screener {
	return screenerStore{
		db: db{pool},
	}
}

// CountScreener implements port.Screener
func (s screenerStore) CountScreener(ctx context.Context, filters domain.ScreenerFilter) (count int, err error) {
	c := Company.Alias("c")

	_, joins := screenerQuery(filters, c)
	cond := screenerFilter(filters)

	row := s.db.QueryRow(ctx, `
			select
				count(distinct c.id)
			from %T
			%T
			where %c
		`, c, pg.Multi(joins), cond)

	err = row.Scan(&count)

	return
}

// IterateScreener implements port.Screener
func (s screenerStore) IterateScreener(ctx context.Context, filters domain.ScreenerFilter) iter.Seq2[*domain.Screener, error] {
	return func(yield func(*domain.Screener, error) bool) {
		c := Company.Alias("c")
		curr := Currency.Alias("curr")

		cols, joins := screenerQuery(filters, c)
		cond := screenerFilter(filters)
		orderBy := screenerOrderBy(filters.OrderBy)

		rows, err := s.db.Query(ctx, `
			select
				%T
			from %T
			left join %T on curr.id = c."currencyId"
			%T
			where %c
			order by %T nulls last
			offset %d
			limit %d
		`, cols, c, curr, pg.Multi(joins), cond, pg.Order(orderBy, filters.Order), filters.Offset, filters.Limit)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var screener domain.Screener

			if err = rows.Scan(
				screenerScanColumns(&screener, filters.Columns)...,
			); err != nil {
				yield(nil, err)
				return
			}

			screenerTransform(&screener)

			if !yield(&screener, nil) {
				return
			}
		}
	}
}

func screenerTransform(screener *domain.Screener) {
	screener.Revenue.Content *= (1_000_000 / 100)
	screener.CostOfRevenue.Content *= (1_000_000 / 100)
	screener.GrossOperatingProfit.Content *= (1_000_000 / 100)
	screener.Ebit.Content *= (1_000_000 / 100)
	screener.NetIncome.Content *= (1_000_000 / 100)
	screener.TotalAssets.Content *= (1_000_000 / 100)
	screener.TotalLiabilities.Content *= (1_000_000 / 100)
	screener.CashAndEquivalents.Content *= (1_000_000 / 100)
	screener.ShortTermInvestments.Content *= (1_000_000 / 100)
	screener.LongTermDebt.Content *= (1_000_000 / 100)
	screener.CurrentDebt.Content *= (1_000_000 / 100)
	screener.Equity.Content *= (1_000_000 / 100)
	screener.OperatingCashFlow.Content *= (1_000_000 / 100)
	screener.CapitalExpenditures.Content *= (1_000_000 / 100)
	screener.FreeCashFlow.Content *= (1_000_000 / 100)
	screener.PPE.Content *= (1_000_000 / 100)
}

func screenerOrderBy(orderBy string) any {
	m := MagicFormulaRankings.Alias("m")
	f := Financials.Alias("f")
	c := Company.Alias("c")
	df := DerivedFinancials.Alias("df")
	sec := Sector.Alias("sec")
	curr := Currency.Alias("curr")

	switch orderBy {
	case domain.ScreenerColumnName:
		return c.Col("name")
	case domain.ScreenerColumnMagicRank:
		return m.Col("rank")
	case domain.ScreenerColumnSector:
		return sec.Col("name")
	case "currency":
		return curr.Col("name")
	// Static financials
	case "revenue", "cost_of_revenue", "gross_operating_profit", "ebit", "net_income", "total_assets", "total_liabilities", "cash_and_equivalents", "short_term_investments", "long_term_debt", "current_debt", "equity", "operating_cash_flow", "capital_expenditures", "free_cash_flow", "number_of_shares", "ppe":
		return f.Col(orderBy)
		// Derived financials
	case "eps", "pe", "evebit", "ps", "pb":
		return df.Col(orderBy)
	default:
		return orderBy
	}
}

func screenerFilter(filters domain.ScreenerFilter) pg.QueryEncoder {
	// m := MagicFormulaRankings.Alias("m")
	f := Financials.Alias("f")
	// sec := Sector.Alias("sec")
	// df := DerivedFinancials.Alias("df")

	cond := pg.And()

	if filters.Revenue.Min.Valid {
		cond.And(pg.Gte(f.Col("revenue"), filters.Revenue.Min))
	}

	if filters.Revenue.Max.Valid {
		cond.And(pg.Lte(f.Col("revenue"), filters.Revenue.Max))
	}

	return cond
}

func screenerQuery(filters domain.ScreenerFilter, a pg.Alias) (columns, []pg.QueryEncoder) {
	curr := Currency.Alias("curr")

	cols := make([]pg.ChainedIdentifier, 4, len(filters.Columns)+4)
	cols[0] = a.Col("id")
	cols[1] = a.Col("name")
	cols[2] = curr.Col("name")
	cols[3] = a.Col("country_code")
	joins := make([]pg.QueryEncoder, 0)
	tables := make(map[pg.Identifier]struct{})

	m := MagicFormulaRankings.Alias("m")
	f := Financials.Alias("f")
	sec := Sector.Alias("sec")
	df := DerivedFinancials.Alias("df")

	for _, c := range filters.Columns {
		switch c {
		case domain.ScreenerColumnMagicRank:
			cols = append(cols, m.Col("rank"))
			if join := tryAddTable(tables, MagicFormulaRankings, a, m, filters); join != nil {
				joins = append(joins, join)
			}

		case domain.ScreenerColumnSector:
			cols = append(cols, sec.Col("name"))

			if _, ok := tables[Sector]; ok {
				continue
			}

			joins = append(joins, pg.Raw("left join %T on %c", sec, pg.Eq(sec.Col("id"), a.Col("sectorId"))))

		// Static financials
		case "revenue", "cost_of_revenue", "gross_operating_profit", "ebit", "net_income", "total_assets", "total_liabilities", "cash_and_equivalents", "short_term_investments", "long_term_debt", "current_debt", "equity", "operating_cash_flow", "capital_expenditures", "free_cash_flow", "number_of_shares", "ppe":
			cols = append(cols, f.Col(c))
			if join := tryAddTable(tables, Financials, a, f, filters); join != nil {
				joins = append(joins, join)
			}

		// Derived financials
		case "eps", "pe", "evebit", "ps", "pb":
			cols = append(cols, df.Col(c))
			if join := tryAddTable(tables, DerivedFinancials, a, df, filters); join != nil {
				joins = append(joins, join)
			}
		}
	}

	return columns{cols}, joins
}

func tryAddTable(tables map[pg.Identifier]struct{}, identifier pg.Identifier, a1 pg.Alias, a2 pg.Alias, filters domain.ScreenerFilter) pg.QueryEncoder {
	if _, ok := tables[identifier]; ok {
		return nil
	}
	tables[identifier] = struct{}{}

	return pg.Raw("left join %T on %c", a2, pg.And(pg.Eq(a2.Col("company_id"), a1.Col("id")), pg.Eq(a2.Col("fiscal_year"), filters.FiscalYear)))
}

func screenerScanColumns(screener *domain.Screener, cols []string) []any {
	scans := make([]any, 4, len(cols)+4)

	scans[0] = &screener.CompanyId
	scans[1] = &screener.Name
	scans[2] = &screener.Currency
	scans[3] = &screener.CountryCode

	//TODO: Consider replacing this with reflection or cache a reflected domain.Screener to set values.
	for _, c := range cols {
		switch c {
		case domain.ScreenerColumnMagicRank:
			scans = append(scans, &screener.MagicRank)
		case domain.ScreenerColumnSector:
			scans = append(scans, &screener.Sector)

		// Static financials
		case domain.ScreenerColumnRevenue:
			scans = append(scans, &screener.Revenue)
		case "cost_of_revenue":
			scans = append(scans, &screener.CostOfRevenue)
		case "gross_operating_profit":
			scans = append(scans, &screener.GrossOperatingProfit)
		case "ebit":
			scans = append(scans, &screener.Ebit)
		case "net_income":
			scans = append(scans, &screener.NetIncome)
		case "total_assets":
			scans = append(scans, &screener.TotalAssets)
		case "total_liabilities":
			scans = append(scans, &screener.TotalLiabilities)
		case "cash_and_equivalents":
			scans = append(scans, &screener.CashAndEquivalents)
		case "short_term_investments":
			scans = append(scans, &screener.ShortTermInvestments)
		case "long_term_debt":
			scans = append(scans, &screener.LongTermDebt)
		case "current_debt":
			scans = append(scans, &screener.CurrentDebt)
		case "equity":
			scans = append(scans, &screener.Equity)
		case "operating_cash_flow":
			scans = append(scans, &screener.OperatingCashFlow)
		case "capital_expenditures":
			scans = append(scans, &screener.CapitalExpenditures)
		case "free_cash_flow":
			scans = append(scans, &screener.FreeCashFlow)
		case "number_of_shares":
			scans = append(scans, &screener.NumberOfShares)
		case "ppe":
			scans = append(scans, &screener.PPE)

		// Derived financials
		case "eps":
			scans = append(scans, &screener.Eps)
		case "pe":
			scans = append(scans, &screener.Pe)
		case "evebit":
			scans = append(scans, &screener.Evebit)
		case "ps":
			scans = append(scans, &screener.Ps)
		case "pb":
			scans = append(scans, &screener.Pb)
		}
	}

	return scans
}

// CountMagicRanks implements port.Screener
func (s screenerStore) CountMagicRanks(ctx context.Context, filters domain.MagicRankFilter) (count int, err error) {
	m := MagicFormulaRankings.Alias("m")

	cond := MagicRanksFilter(filters, m)

	row := s.db.QueryRow(ctx, `
			select
				count(*)
			from %T
			where %c
		`, m, cond)

	err = row.Scan(&count)

	return
}

// IterateMagicRanks implements port.Screener
func (s screenerStore) IterateMagicRanks(ctx context.Context, filters domain.MagicRankFilter) iter.Seq2[*domain.MagicRank, error] {
	return func(yield func(*domain.MagicRank, error) bool) {
		m := MagicFormulaRankings.Alias("m")
		cond := pg.And()

		if filters.FiscalYear != 0 {
			cond.And(pg.Eq(m.Col("fiscal_year"), filters.FiscalYear))
		}

		rows, err := s.db.Query(ctx, `
			select
				m.company_id,
				c.name,
				m.fiscal_year,
				m.roc,
				m.yield,
				m.roc_rank,
				m.yield_rank,
				m.rank
			from %T
			left join %T c on c.id = m.company_id
			WHERE %c
		`, m, Company, cond)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var rank domain.MagicRank

			if err = rows.Scan(
				&rank.Company.ID,
				&rank.Company.Name,
				&rank.FiscalYear,
				&rank.ROC,
				&rank.EarningsYield,
				&rank.ROCRank,
				&rank.EarningsYieldRank,
				&rank.Rank,
			); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&rank, nil) {
				return
			}
		}
	}
}

func MagicRanksFilter(filters domain.MagicRankFilter, a pg.Alias) pg.QueryEncoder {
	cond := pg.And()

	return cond
}
