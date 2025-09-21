package domain

import "github.com/rs/xid"

const (
	TaxRate = 20.6
)

type Financials struct {
	CompanyID   xid.ID               `json:"companyId"`
	FiscalYear  int                  `json:"fiscalYear"`
	CurrencyID  xid.ID               `json:"currency"`
	StaticData  FinancialData        `json:"staticData"`
	DerivedData DerivedFinancialData `json:"derivedData"`
}

type FinancialData struct {
	CapitalExpenditures  int `json:"capital_expenditures"`
	CashAndEquivalents   int `json:"cash_and_equivalents"`
	CostOfRevenue        int `json:"cost_of_revenue"`
	CurrentDebt          int `json:"current_debt"`
	Ebit                 int `json:"ebit"`
	Equity               int `json:"equity"`
	FreeCashFlow         int `json:"free_cash_flow"`
	GrossOperatingProfit int `json:"gross_operating_profit"`
	LongTermDebt         int `json:"long_term_debt"`
	NetIncome            int `json:"net_income"`
	NumberOfShares       int `json:"number_of_shares"`
	OperatingCashFlow    int `json:"operating_cash_flow"`
	PPE                  int `json:"ppe"`
	Revenue              int `json:"revenue"`
	ShortTermInvestments int `json:"short_term_investments"`
	TotalAssets          int `json:"total_assets"`
	TotalLiabilities     int `json:"total_liabilities"`
}

// Derived financials
type DerivedFinancialData struct {
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

type FinancialFilter struct {
	Order   string   `query:"order" enum:"asc,desc" default:"asc"`
	OrderBy string   `query:"orderBy" enum:"name" default:"name"`
	Limit   int      `query:"limit" min:"1" max:"500" default:"50"`
	Offset  int      `query:"offset" min:"0"`
	Include []xid.ID `query:"include"`
	Search  string   `query:"search"`
}
