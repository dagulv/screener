package domain

import (
	"github.com/rs/xid"
)

type Screener struct {
	CompanyId   xid.ID           `json:"companyId"`
	Name        string           `json:"name"`
	Currency    string           `json:"currency"`
	CountryCode CountryCode      `json:"countryCode"`
	MagicRank   Nullable[int64]  `json:"magicRank"`
	Sector      Nullable[string] `json:"sector"`

	// Static financials
	CapitalExpenditures  Nullable[int64] `json:"capital_expenditures"`
	CashAndEquivalents   Nullable[int64] `json:"cash_and_equivalents"`
	CostOfRevenue        Nullable[int64] `json:"cost_of_revenue"`
	CurrentDebt          Nullable[int64] `json:"current_debt"`
	Ebit                 Nullable[int64] `json:"ebit"`
	Equity               Nullable[int64] `json:"equity"`
	FreeCashFlow         Nullable[int64] `json:"free_cash_flow"`
	GrossOperatingProfit Nullable[int64] `json:"gross_operating_profit"`
	LongTermDebt         Nullable[int64] `json:"long_term_debt"`
	NetIncome            Nullable[int64] `json:"net_income"`
	NumberOfShares       Nullable[int64] `json:"number_of_shares"`
	OperatingCashFlow    Nullable[int64] `json:"operating_cash_flow"`
	PPE                  Nullable[int64] `json:"ppe"`
	Revenue              Nullable[int64] `json:"revenue"`
	ShortTermInvestments Nullable[int64] `json:"short_term_investments"`
	TotalAssets          Nullable[int64] `json:"total_assets"`
	TotalLiabilities     Nullable[int64] `json:"total_liabilities"`

	// Derived financials
	Eps    Nullable[float64] `json:"eps"`
	Evebit Nullable[float64] `json:"evebit"`
	Pb     Nullable[float64] `json:"pb"`
	Pe     Nullable[float64] `json:"pe"`
	Ps     Nullable[float64] `json:"ps"`
}

type ScreenerFilter struct {
	Order      string      `query:"order" enum:"asc,desc" default:"asc"`
	OrderBy    string      `query:"orderby" enum:"name,magicRank,sector,revenue,cost_of_revenue,gross_operating_profit,ebit,net_income,total_assets,total_liabilities,cash_and_equivalents,short_term_investments,long_term_debt,current_debt,equity,operating_cash_flow,capital_expenditures,free_cash_flow,number_of_shares,ppe,eps,pe,evebit,ps,pb" default:"name"`
	Limit      int         `query:"limit" min:"1" max:"500" default:"50"`
	Offset     int         `query:"offset" min:"0"`
	Include    []xid.ID    `query:"include"`
	Search     string      `query:"search"`
	FiscalYear int         `query:"fiscalYear"`
	Columns    []string    `query:"columns"`
	Revenue    MinMax[int] `query:"revenue"` //min:"0" max:"1000000"
}

type ScreenerColumn struct {
	ID string `query:"id" enum:"magicRank,revenue"`
}

const (
	ScreenerColumnMagicRank = "magicRank"
	ScreenerColumnRevenue   = "revenue"
	ScreenerColumnSector    = "sector"
	ScreenerColumnName      = "name"
)
