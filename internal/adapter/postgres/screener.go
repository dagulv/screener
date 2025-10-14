package postgres

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/webmafia/pg"
)

const (
	FloatConstant     = 100
	TransformConstant = 1_000_000 / FloatConstant
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
		// cr := QuarterlyCurrencyRates.Alias("cr")

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
	screener.Revenue.Content *= TransformConstant
	screener.CostOfRevenue.Content *= TransformConstant
	screener.GrossOperatingProfit.Content *= TransformConstant
	screener.EBIT.Content *= TransformConstant
	screener.NetIncome.Content *= TransformConstant
	screener.TotalAssets.Content *= TransformConstant
	screener.TotalLiabilities.Content *= TransformConstant
	screener.CashAndEquivalents.Content *= TransformConstant
	screener.ShortTermInvestments.Content *= TransformConstant
	screener.LongTermDebt.Content *= TransformConstant
	screener.CurrentDebt.Content *= TransformConstant
	screener.Equity.Content *= TransformConstant
	screener.OperatingCashFlow.Content *= TransformConstant
	screener.CapitalExpenditures.Content *= TransformConstant
	screener.FreeCashFlow.Content *= TransformConstant
	screener.PPE.Content *= TransformConstant
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
	// Static financials with currency
	case "revenue", "cost_of_revenue", "gross_operating_profit", "ebit", "net_income", "total_assets", "total_liabilities", "cash_and_equivalents", "short_term_investments", "long_term_debt", "current_debt", "equity", "operating_cash_flow", "capital_expenditures", "free_cash_flow", "ppe":
		return pg.Col("f." + orderBy + "::real / cr.rate")
	// Static financials without currency
	case "number_of_shares":
		return f.Col(orderBy)
	// Derived financials
	case "eps", "pe", "evebit", "ps", "pb", "operating_margin", "net_margin", "roe", "roc", "liabilities_to_equity", "debt_to_ebit", "debt_to_assets", "cash_conversion":
		return df.Col(orderBy)
	default:
		return orderBy
	}
}

func screenerFilter(filters domain.ScreenerFilter) pg.QueryEncoder {
	c := Company.Alias("c")
	// cr := QuarterlyCurrencyRates.Alias("cr")
	cond := pg.And()
	if filters.Search != "" {
		cond.And(pg.Search(c.Col("ts"), filters.Search, pg.SearchOptions{
			Preprocessor: pg.PrefixSearch,
		}))
		return cond
	}
	m := MagicFormulaRankings.Alias("m")
	// f := Financials.Alias("f")
	// sec := Sector.Alias("sec")
	df := DerivedFinancials.Alias("df")

	if filters.CapitalExpenditures.Min.Valid {
		cond.And(pg.Raw("(f.capital_expenditures::real / cr.rate) >= %c", filters.CapitalExpenditures.Min.Content*FloatConstant))
	}
	if filters.CapitalExpenditures.Max.Valid {
		cond.And(pg.Raw("(f.capital_expenditures::real / cr.rate) <= %c", filters.CapitalExpenditures.Max.Content*FloatConstant))
	}

	if filters.EBIT.Min.Valid {
		cond.And(pg.Raw("(f.ebit::real / cr.rate) >= %c", filters.EBIT.Min.Content*FloatConstant))
	}
	if filters.EBIT.Max.Valid {
		cond.And(pg.Raw("(f.ebit::real / cr.rate) <= %c", filters.EBIT.Max.Content*FloatConstant))
	}

	if filters.Equity.Min.Valid {
		cond.And(pg.Raw("(f.equity::real / cr.rate) >= %c", filters.Equity.Min.Content*FloatConstant))
	}
	if filters.Equity.Max.Valid {
		cond.And(pg.Raw("(f.equity::real / cr.rate) <= %c", filters.Equity.Max.Content*FloatConstant))
	}

	if filters.GrossOperatingProfit.Min.Valid {
		cond.And(pg.Raw("(f.gross_operating_profit::real / cr.rate) >= %c", filters.GrossOperatingProfit.Min.Content*FloatConstant))
	}
	if filters.GrossOperatingProfit.Max.Valid {
		cond.And(pg.Raw("(f.gross_operating_profit::real / cr.rate) <= %c", filters.GrossOperatingProfit.Max.Content*FloatConstant))
	}

	if filters.NetIncome.Min.Valid {
		cond.And(pg.Raw("(f.net_income::real / cr.rate) >= %c", filters.NetIncome.Min.Content*FloatConstant))
	}
	if filters.NetIncome.Max.Valid {
		cond.And(pg.Raw("(f.net_income::real / cr.rate) <= %c", filters.NetIncome.Max.Content*FloatConstant))
	}

	if filters.OperatingCashFlow.Min.Valid {
		cond.And(pg.Raw("(f.operating_cash_flow::real / cr.rate) >= %c", filters.OperatingCashFlow.Min.Content*FloatConstant))
	}
	if filters.OperatingCashFlow.Max.Valid {
		cond.And(pg.Raw("(f.operating_cash_flow::real / cr.rate) <= %c", filters.OperatingCashFlow.Max.Content*FloatConstant))
	}

	if filters.Revenue.Min.Valid {
		// cond.And(pg.Gte(f.Col("revenue"), filters.Revenue.Min.Content*FloatConstant))
		cond.And(pg.Raw("(f.revenue::real / cr.rate) >= %c", filters.Revenue.Min.Content*FloatConstant))
	}
	if filters.Revenue.Max.Valid {
		cond.And(pg.Raw("(f.revenue::real / cr.rate) <= %c", filters.Revenue.Max.Content*FloatConstant))
	}

	if filters.EPS.Min.Valid {
		cond.And(pg.Gte(df.Col("eps"), filters.EPS.Min))
	}
	if filters.EPS.Max.Valid {
		cond.And(pg.Lte(df.Col("eps"), filters.EPS.Max))
	}

	if filters.EVEBIT.Min.Valid {
		cond.And(pg.Gte(df.Col("evebit"), filters.EVEBIT.Min))
	}
	if filters.EVEBIT.Max.Valid {
		cond.And(pg.Lte(df.Col("evebit"), filters.EVEBIT.Max))
	}

	if filters.PB.Min.Valid {
		cond.And(pg.Gte(df.Col("pb"), filters.PB.Min))
	}
	if filters.PB.Max.Valid {
		cond.And(pg.Lte(df.Col("pb"), filters.PB.Max))
	}

	if filters.PE.Min.Valid {
		cond.And(pg.Gte(df.Col("pe"), filters.PE.Min))
	}
	if filters.PE.Max.Valid {
		cond.And(pg.Lte(df.Col("pe"), filters.PE.Max))
	}

	if filters.PS.Min.Valid {
		cond.And(pg.Gte(df.Col("ps"), filters.PS.Min))
	}
	if filters.PS.Max.Valid {
		cond.And(pg.Lte(df.Col("ps"), filters.PS.Max))
	}

	if filters.OperatingMargin.Min.Valid {
		cond.And(pg.Gte(df.Col("operating_margin"), filters.OperatingMargin.Min))
	}
	if filters.OperatingMargin.Max.Valid {
		cond.And(pg.Lte(df.Col("operating_margin"), filters.OperatingMargin.Max))
	}

	if filters.NetMargin.Min.Valid {
		cond.And(pg.Gte(df.Col("net_margin"), filters.NetMargin.Min))
	}
	if filters.NetMargin.Max.Valid {
		cond.And(pg.Lte(df.Col("net_margin"), filters.NetMargin.Max))
	}

	if filters.ROE.Min.Valid {
		cond.And(pg.Gte(df.Col("roe"), filters.ROE.Min))
	}
	if filters.ROE.Max.Valid {
		cond.And(pg.Lte(df.Col("roe"), filters.ROE.Max))
	}

	if filters.ROC.Min.Valid {
		cond.And(pg.Gte(df.Col("roc"), filters.ROC.Min))
	}
	if filters.ROC.Max.Valid {
		cond.And(pg.Lte(df.Col("roc"), filters.ROC.Max))
	}

	if filters.LiabilitiesToEquity.Min.Valid {
		cond.And(pg.Gte(df.Col("liabilities_to_equity"), filters.LiabilitiesToEquity.Min))
	}
	if filters.LiabilitiesToEquity.Max.Valid {
		cond.And(pg.Lte(df.Col("liabilities_to_equity"), filters.LiabilitiesToEquity.Max))
	}

	if filters.DebtToEBIT.Min.Valid {
		cond.And(pg.Gte(df.Col("debt_to_ebit"), filters.DebtToEBIT.Min))
	}
	if filters.DebtToEBIT.Max.Valid {
		cond.And(pg.Lte(df.Col("debt_to_ebit"), filters.DebtToEBIT.Max))
	}

	if filters.DebtToAssets.Min.Valid {
		cond.And(pg.Gte(df.Col("debt_to_assets"), filters.DebtToAssets.Min))
	}
	if filters.DebtToAssets.Max.Valid {
		cond.And(pg.Lte(df.Col("debt_to_assets"), filters.DebtToAssets.Max))
	}

	if filters.CashConversion.Min.Valid {
		cond.And(pg.Gte(df.Col("cash_conversion"), filters.CashConversion.Min))
	}
	if filters.CashConversion.Max.Valid {
		cond.And(pg.Lte(df.Col("cash_conversion"), filters.CashConversion.Max))
	}

	if filters.MagicRank.Min.Valid {
		cond.And(pg.Gte(m.Col("rank"), filters.MagicRank.Min))
	}
	if filters.MagicRank.Max.Valid {
		cond.And(pg.Lte(m.Col("rank"), filters.MagicRank.Max))
	}

	return cond
}

func screenerQuery(filters domain.ScreenerFilter, a pg.Alias) (columns, []pg.QueryEncoder) {
	curr := Currency.Alias("curr")
	cr := QuarterlyCurrencyRates.Alias("cr")

	cols := make([]pg.ChainedIdentifier, 4, len(filters.Columns)+4)
	cols[0] = a.Col("id")
	cols[1] = a.Col("name")
	cols[2] = curr.Col("name")
	cols[3] = a.Col("country_code")
	joins := make([]pg.QueryEncoder, 1)
	joins[0] = pg.Raw("left join %T on %c", cr, pg.And(pg.Eq(cr.Col("fiscal_year"), filters.FiscalYear), pg.Eq(cr.Col("currency_id"), a.Col("currencyId")), pg.Eq(cr.Col("quarter"), 1)))
	// joins[0] = pg.Raw(`left join quarterly_currency_rates cr on cr.fiscal_year = 2024 and cr.currency_id = c."currencyId" and cr.quarter = 1`)
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
		case "eps", "pe", "evebit", "ps", "pb", "operating_margin", "net_margin", "roe", "roc", "liabilities_to_equity", "debt_to_ebit", "debt_to_assets", "cash_conversion":
			cols = append(cols, df.Col(c))
			if join := tryAddTable(tables, DerivedFinancials, a, df, filters); join != nil {
				joins = append(joins, join)
			}
		}
	}

	if domain.HasValid(filters.CapitalExpenditures, filters.EBIT, filters.Equity, filters.GrossOperatingProfit, filters.NetIncome, filters.OperatingCashFlow, filters.Revenue) {
		if join := tryAddTable(tables, Financials, a, f, filters); join != nil {
			joins = append(joins, join)
		}
	}

	if domain.HasValid(filters.EPS,
		filters.EVEBIT,
		filters.PB,
		filters.PE,
		filters.PS,
		filters.OperatingMargin,
		filters.NetMargin,
		filters.ROE,
		filters.ROC,
		filters.LiabilitiesToEquity,
		filters.DebtToEBIT,
		filters.DebtToAssets,
		filters.CashConversion,
		filters.MagicRank) {
		if join := tryAddTable(tables, DerivedFinancials, a, df, filters); join != nil {
			joins = append(joins, join)
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
			scans = append(scans, &screener.EBIT)
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
			scans = append(scans, &screener.EPS)
		case "pe":
			scans = append(scans, &screener.PE)
		case "evebit":
			scans = append(scans, &screener.EVEBIT)
		case "ps":
			scans = append(scans, &screener.PS)
		case "pb":
			scans = append(scans, &screener.PB)
		case "operating_margin":
			scans = append(scans, &screener.OperatingMargin)
		case "net_margin":
			scans = append(scans, &screener.NetMargin)
		case "roe":
			scans = append(scans, &screener.ROE)
		case "roc":
			scans = append(scans, &screener.ROC)
		case "liabilities_to_equity":
			scans = append(scans, &screener.LiabilitiesToEquity)
		case "debt_to_ebit":
			scans = append(scans, &screener.DebtToEbit)
		case "debt_to_assets":
			scans = append(scans, &screener.DebtToAssets)
		case "cash_conversion":
			scans = append(scans, &screener.CashConversion)
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
