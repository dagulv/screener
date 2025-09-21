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
	EBIT                 Nullable[int64] `json:"ebit"`
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
	EPS                 Nullable[float64] `json:"eps"`
	EVEBIT              Nullable[float64] `json:"evebit"`
	PB                  Nullable[float64] `json:"pb"`
	PE                  Nullable[float64] `json:"pe"`
	PS                  Nullable[float64] `json:"ps"`
	OperatingMargin     Nullable[float64] `json:"operating_margin"`
	NetMargin           Nullable[float64] `json:"net_margin"`
	ROE                 Nullable[float64] `json:"roe"`
	ROC                 Nullable[float64] `json:"roc"`
	LiabilitiesToEquity Nullable[float64] `json:"liabilities_to_equity"`
	DebtToEbit          Nullable[float64] `json:"debt_to_ebit"`
	DebtToAssets        Nullable[float64] `json:"debt_to_assets"`
	CashConversion      Nullable[float64] `json:"cash_conversion"`
}

type ScreenerFilter struct {
	Order      string   `query:"order" enum:"asc,desc" default:"asc"`
	OrderBy    string   `query:"orderby" enum:"name,magicRank,sector,revenue,cost_of_revenue,gross_operating_profit,ebit,net_income,total_assets,total_liabilities,cash_and_equivalents,short_term_investments,long_term_debt,current_debt,equity,operating_cash_flow,capital_expenditures,free_cash_flow,number_of_shares,ppe,eps,pe,evebit,ps,pb" default:"name"`
	Limit      int      `query:"limit" min:"1" max:"500" default:"50"`
	Offset     int      `query:"offset" min:"0"`
	Include    []xid.ID `query:"include"`
	Search     string   `query:"search"`
	FiscalYear int      `query:"fiscalYear"`
	Columns    []string `query:"columns"`

	// Static financials
	CapitalExpenditures  MinMax[int] `query:"capital_expenditures"`   // min="0" max="1000000"
	EBIT                 MinMax[int] `query:"ebit"`                   // min="0" max="1000000"
	Equity               MinMax[int] `query:"equity"`                 // min="0" max="1000000"
	GrossOperatingProfit MinMax[int] `query:"gross_operating_profit"` // min="0" max="1000000"
	NetIncome            MinMax[int] `query:"net_income"`             // min="0" max="1000000"
	OperatingCashFlow    MinMax[int] `query:"operating_cash_flow"`    // min="0" max="1000000"
	Revenue              MinMax[int] `query:"revenue"`                // min="0" max="1000000"

	// Derived financials
	EPS                 MinMax[float32] `query:"eps"`                   // min="0" max="1000"
	EVEBIT              MinMax[float32] `query:"evebit"`                // min="0" max="12"
	PB                  MinMax[float32] `query:"pb"`                    // min="0" max="3"
	PE                  MinMax[float32] `query:"pe"`                    // min="5" max="25"
	PS                  MinMax[float32] `query:"ps"`                    // min="0" max="5"
	OperatingMargin     MinMax[float32] `query:"operating_margin"`      // min="0" max="1"
	NetMargin           MinMax[float32] `query:"net_margin"`            // min="0" max="1"
	ROE                 MinMax[float32] `query:"roe"`                   // min="0" max="1"
	ROC                 MinMax[float32] `query:"roc"`                   // min="0" max="1"
	LiabilitiesToEquity MinMax[float32] `query:"liabilities_to_equity"` // min="0" max="2"
	DebtToEBIT          MinMax[float32] `query:"debt_to_ebit"`          // min="0" max="3"
	DebtToAssets        MinMax[float32] `query:"debt_to_assets"`        // min="0" max="1"
	CashConversion      MinMax[float32] `query:"cash_conversion"`       // min="0" max="2"
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
