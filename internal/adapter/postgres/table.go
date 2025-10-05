package postgres

import "github.com/webmafia/pg"

const (
	Sector                 pg.Identifier = "sectors"
	Currency               pg.Identifier = "currencies"
	Company                pg.Identifier = "companies"
	Financials             pg.Identifier = "financials"
	DerivedFinancials      pg.Identifier = "derived_financials"
	Share                  pg.Identifier = "shares"
	MagicFormulaRankings   pg.Identifier = "magic_formula_rankings"
	QuarterlyCurrencyRates pg.Identifier = "quarterly_currency_rates"
)
