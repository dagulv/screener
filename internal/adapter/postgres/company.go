package postgres

import (
	"context"
	"iter"

	"github.com/dagulv/screener/internal/core/domain"
	"github.com/dagulv/screener/internal/core/port"
	"github.com/rs/xid"
	"github.com/webmafia/pg"
)

type companyStore struct {
	db
}

func NewCompany(pool *pg.DB) port.Company {
	return companyStore{
		db: db{pool},
	}
}

// CreateCompany implements port.Company
func (s companyStore) CreateCompany(ctx context.Context, company *domain.Company) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("id", company.ID).
		Value("name", company.Name).
		Value("symbol", company.Symbol).
		Value("isin", company.ISIN).
		Value("currencyId", company.Currency.ID).
		Value("sectorId", company.Sector.ID).
		Value("orderbookId", company.OrderbookID).
		Value("country_code", company.CountryCode).
		Value("market_place_code", company.MarketPlaceCode)

	_, err = s.db.InsertValues(ctx, Company, vals)

	return
}

// ReadCompany implements port.Company
func (s companyStore) ReadCompany(ctx context.Context, company *domain.Company) (err error) {
	c := Company.Alias("c")

	row := s.db.QueryRow(ctx, `
		select
			c.id,
			c.name,
			c.symbol,
			c.isin,
			curr.id,
			curr.name,
			sec.id,
			sec.name,
			c."orderbookId",
			c.country_code,
			c.market_place_code
		from %T
		left join %T curr on curr.id = c."currencyId"
		left join %T sec on sec.id = c."sectorId"
		where %c
	`, c, Currency, Sector, pg.Eq(c.Col("id"), company.ID))

	err = row.Scan(
		&company.ID,
		&company.Name,
		&company.Symbol,
		&company.ISIN,
		&company.Currency.ID,
		&company.Currency.Name,
		&company.Sector.ID,
		&company.Sector.Name,
		&company.OrderbookID,
		&company.CountryCode,
		&company.MarketPlaceCode,
	)

	return
}

// UpdateCompany implements port.Company
func (s companyStore) UpdateCompany(ctx context.Context, company *domain.Company) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("name", company.Name).
		Value("symbol", company.Symbol).
		Value("isin", company.ISIN).
		Value("currencyId", company.Currency.ID).
		Value("sectorId", company.Sector.ID).
		Value("orderbookId", company.OrderbookID).
		Value("country_code", company.CountryCode).
		Value("market_place_code", company.MarketPlaceCode)

	_, err = s.db.UpdateValues(ctx, Company, vals, pg.Eq("id", company.ID))

	return
}

// DeleteCompany implements port.Company
func (s companyStore) DeleteCompany(ctx context.Context, companyId xid.ID) (err error) {
	_, err = s.db.Delete(ctx, Company, pg.Eq("id", companyId))
	return
}

// CountCompanies implements port.Company
func (s companyStore) CountCompanies(ctx context.Context, filters domain.CompanyFilter) (count int, err error) {
	c := Company.Alias("c")

	cond := companiesFilter(filters, c)

	row := s.db.QueryRow(ctx, `
			select
				count(*)
			from %T
			where %c
		`, c, cond)

	err = row.Scan(&count)

	return
}

// IterateCompanies implements port.Company
func (s companyStore) IterateCompanies(ctx context.Context, filters domain.CompanyFilter) iter.Seq2[*domain.Company, error] {
	return func(yield func(*domain.Company, error) bool) {
		c := Company.Alias("c")
		curr := Currency.Alias("curr")
		sec := Sector.Alias("sec")

		cond := companiesFilter(filters, c)

		rows, err := s.db.Query(ctx, `
			select
				c.id,
				c.name,
				c.symbol,
				c.isin,
				curr.id,
				curr.name,
				sec.id,
				sec.name,
				c."orderbookId",
				c.country_code,
				c.market_place_code
			from %T
			left join %T on curr.id = c."currencyId"
			left join %T on sec.id = c."sectorId"
			where %c
			order by %T
			offset %T
			limit %T
		`, c, curr, sec, cond, pg.Order(c.Col(filters.OrderBy), filters.Order), filters.Offset, filters.Limit)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var company domain.Company

			if err = rows.Scan(
				&company.ID,
				&company.Name,
				&company.Symbol,
				&company.ISIN,
				&company.Currency.ID,
				&company.Currency.Name,
				&company.Sector.ID,
				&company.Sector.Name,
				&company.OrderbookID,
				&company.CountryCode,
				&company.MarketPlaceCode,
			); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&company, nil) {
				return
			}
		}
	}
}

func companiesFilter(filters domain.CompanyFilter, c pg.Alias) pg.QueryEncoder {
	cond := pg.And()

	if filters.Search != "" {
		cond.And(pg.Raw(`c.name ilike %c`, "%"+filters.Search+"%"))
	}

	return cond
}

// CreateFinancials implements port.Company
func (s companyStore) CreateFinancials(ctx context.Context, financials *domain.Financials) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("company_id", financials.CompanyID).
		Value("fiscal_year", financials.FiscalYear).
		Value("currency", financials.CurrencyID).
		Value("revenue", financials.StaticData.Revenue).
		Value("cost_of_revenue", financials.StaticData.CostOfRevenue).
		Value("gross_operating_profit", financials.StaticData.GrossOperatingProfit).
		Value("ebit", financials.StaticData.Ebit).
		Value("net_income", financials.StaticData.NetIncome).
		Value("total_assets", financials.StaticData.TotalAssets).
		Value("total_liabilities", financials.StaticData.TotalLiabilities).
		Value("cash_and_equivalents", financials.StaticData.CashAndEquivalents).
		Value("short_term_investments", financials.StaticData.ShortTermInvestments).
		Value("long_term_debt", financials.StaticData.LongTermDebt).
		Value("current_debt", financials.StaticData.CurrentDebt).
		Value("equity", financials.StaticData.Equity).
		Value("operating_cash_flow", financials.StaticData.OperatingCashFlow).
		Value("capital_expenditures", financials.StaticData.CapitalExpenditures).
		Value("free_cash_flow", financials.StaticData.FreeCashFlow).
		Value("number_of_shares", financials.StaticData.NumberOfShares).
		Value("ppe", financials.StaticData.PPE)

	// _, err = s.db.InsertValues(ctx, Financials, vals, pg.InsertOptions{OnConflict: pg.DoUpdate(2, "company_id", "fiscal_year")})
	_, err = s.db.UpdateValues(ctx, Financials, vals, pg.And(pg.Eq("company_id", financials.CompanyID), pg.Eq("fiscal_year", financials.FiscalYear)))

	return
}

// CountFinancials implements port.Company
func (s companyStore) CountFinancials(ctx context.Context, filters domain.FinancialFilter) (count int, err error) {
	f := Financials.Alias("f")
	cond := financialsFilter(filters, f)

	row := s.db.QueryRow(ctx, `
		select
			count(*)
		from %T
		where %c
	`, f, cond)

	err = row.Scan(&count)

	return
}

// IterateFinancials implements port.Company
func (s companyStore) IterateFinancials(ctx context.Context, filters domain.FinancialFilter) iter.Seq2[*domain.Financials, error] {
	return func(yield func(*domain.Financials, error) bool) {
		f := Financials.Alias("f")
		df := DerivedFinancials.Alias("df")
		cond := financialsFilter(filters, f)

		rows, err := s.db.Query(ctx, `
			select
				f.company_id,
				f.fiscal_year,
				f.currency,
				f.revenue,
				f.cost_of_revenue,
				f.gross_operating_profit,
				f.ebit,
				f.net_income,
				f.total_assets,
				f.total_liabilities,
				f.cash_and_equivalents,
				f.short_term_investments,
				f.long_term_debt,
				f.current_debt,
				f.equity,
				f.operating_cash_flow,
				f.capital_expenditures,
				f.free_cash_flow,
				f.number_of_shares,
				f.ppe,
				df.eps,
				df.evebit,
				df.pb,
				df.pe,
				df.ps
			from %T
			left join %T on df.company_id = f.company_id and df.fiscal_year = f.fiscal_year
			where %c
		`, f, df, cond)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var financials domain.Financials

			if err = rows.Scan(
				&financials.CompanyID,
				&financials.FiscalYear,
				&financials.CurrencyID,
				&financials.StaticData.Revenue,
				&financials.StaticData.CostOfRevenue,
				&financials.StaticData.GrossOperatingProfit,
				&financials.StaticData.Ebit,
				&financials.StaticData.NetIncome,
				&financials.StaticData.TotalAssets,
				&financials.StaticData.TotalLiabilities,
				&financials.StaticData.CashAndEquivalents,
				&financials.StaticData.ShortTermInvestments,
				&financials.StaticData.LongTermDebt,
				&financials.StaticData.CurrentDebt,
				&financials.StaticData.Equity,
				&financials.StaticData.OperatingCashFlow,
				&financials.StaticData.CapitalExpenditures,
				&financials.StaticData.FreeCashFlow,
				&financials.StaticData.NumberOfShares,
				&financials.StaticData.PPE,
				&financials.DerivedData.EPS,
				&financials.DerivedData.EVEBIT,
				&financials.DerivedData.PB,
				&financials.DerivedData.PE,
				&financials.DerivedData.PS,
			); err != nil {
				yield(nil, err)
				return
			}

			financialsTransform(&financials.StaticData)

			if !yield(&financials, nil) {
				return
			}
		}
	}
}

func financialsTransform(financials *domain.FinancialData) {
	financials.Revenue *= (1_000_000 / 100)
	financials.CostOfRevenue *= (1_000_000 / 100)
	financials.GrossOperatingProfit *= (1_000_000 / 100)
	financials.Ebit *= (1_000_000 / 100)
	financials.NetIncome *= (1_000_000 / 100)
	financials.TotalAssets *= (1_000_000 / 100)
	financials.TotalLiabilities *= (1_000_000 / 100)
	financials.CashAndEquivalents *= (1_000_000 / 100)
	financials.ShortTermInvestments *= (1_000_000 / 100)
	financials.LongTermDebt *= (1_000_000 / 100)
	financials.CurrentDebt *= (1_000_000 / 100)
	financials.Equity *= (1_000_000 / 100)
	financials.OperatingCashFlow *= (1_000_000 / 100)
	financials.CapitalExpenditures *= (1_000_000 / 100)
	financials.FreeCashFlow *= (1_000_000 / 100)
	financials.PPE *= (1_000_000 / 100)
}

func financialsFilter(filters domain.FinancialFilter, a pg.Alias) pg.QueryEncoder {
	cond := pg.And()

	if len(filters.Include) > 0 {
		cond.And(pg.In(a.Col("company_id"), filters.Include))
	}

	return cond
}

// CreateShare implements port.Company
func (s companyStore) CreateShare(ctx context.Context, share *domain.Share) (err error) {
	vals := s.db.AcquireValues()
	defer s.db.ReleaseValues(vals)

	vals.
		Value("company_id", share.CompanyID).
		Value("date", share.Date).
		Value("open", share.Open).
		Value("high", share.High).
		Value("low", share.Low).
		Value("close", share.Close).
		Value("volume", share.Volume).
		Value("average", share.Average)

	_, err = s.db.InsertValues(ctx, Share, vals, pg.InsertOptions{OnConflict: pg.DoUpdate(2, "company_id", "date")})

	return
}

// IterateShares implements port.Company
func (s companyStore) IterateShares(ctx context.Context, filters domain.CompanyFilter) iter.Seq2[*domain.Share, error] {
	return func(yield func(*domain.Share, error) bool) {
		sh := Share.Alias("s")
		cond := pg.And()

		rows, err := s.db.Query(ctx, `
			select
				s.company_id,
				s.date,
				s.open::float,
				s.high::float,
				s.low::float,
				s.close::float,
				s.volume,
				s.average::float
			from %T
			where %c
		`, sh, cond)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var share domain.Share

			if err = rows.Scan(
				&share.CompanyID,
				&share.Date,
				&share.Open,
				&share.High,
				&share.Low,
				&share.Close,
				&share.Volume,
				&share.Average,
			); err != nil {
				yield(nil, err)
				return
			}

			shareTransform(&share)

			if !yield(&share, nil) {
				return
			}
		}
	}
}

func shareTransform(share *domain.Share) {
	share.Open /= 100
	share.High /= 100
	share.Low /= 100
	share.Close /= 100
	share.Average /= 100
}

// IterateFinancials implements port.Company
func (s companyStore) IterateFinancialsByMissingShare(ctx context.Context) iter.Seq2[*domain.Financials, error] {
	return func(yield func(*domain.Financials, error) bool) {
		f := Financials.Alias("f")

		rows, err := s.db.Query(ctx, `
			select
				f.company_id,
				f.fiscal_year,
				f.currency,
				f.revenue,
				f.cost_of_revenue,
				f.gross_operating_profit,
				f.ebit,
				f.net_income,
				f.total_assets,
				f.total_liabilities,
				f.cash_and_equivalents,
				f.short_term_investments,
				f.long_term_debt,
				f.current_debt,
				f.equity,
				f.operating_cash_flow,
				f.capital_expenditures,
				f.free_cash_flow
			from %T
			left join shares s
				ON f.company_id = s.company_id
				and s.date = to_date((f.fiscal_year + 1) || '-01-02', 'YYYY-MM-DD')
			where s.date is null;
		`, f)

		if err != nil {
			yield(nil, err)
			return
		}

		defer rows.Close()

		for rows.Next() {
			var financials domain.Financials

			if err = rows.Scan(
				&financials.CompanyID,
				&financials.FiscalYear,
				&financials.CurrencyID,
				&financials.StaticData.Revenue,
				&financials.StaticData.CostOfRevenue,
				&financials.StaticData.GrossOperatingProfit,
				&financials.StaticData.Ebit,
				&financials.StaticData.NetIncome,
				&financials.StaticData.TotalAssets,
				&financials.StaticData.TotalLiabilities,
				&financials.StaticData.CashAndEquivalents,
				&financials.StaticData.ShortTermInvestments,
				&financials.StaticData.LongTermDebt,
				&financials.StaticData.CurrentDebt,
				&financials.StaticData.Equity,
				&financials.StaticData.OperatingCashFlow,
				&financials.StaticData.CapitalExpenditures,
				&financials.StaticData.FreeCashFlow,
			); err != nil {
				yield(nil, err)
				return
			}

			if !yield(&financials, nil) {
				return
			}
		}
	}
}
