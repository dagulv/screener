// Static financials
export const staticFinancials: Record<
	string,
	{ label: string; style: Intl.NumberFormatOptions['style'] }
> = {
	revenue: { label: 'Revenue', style: 'currency' },
	cost_of_revenue: { label: 'Cost of revenue', style: 'currency' },
	gross_operating_profit: { label: 'Gross operating profit', style: 'currency' },
	ebit: { label: 'EBIT', style: 'currency' },
	net_income: { label: 'Net income', style: 'currency' },
	total_assets: { label: 'Total assets', style: 'currency' },
	total_liabilities: { label: 'Total liabilities', style: 'currency' },
	cash_and_equivalents: { label: 'Cash and equivalents', style: 'currency' },
	short_term_investments: { label: 'Short term investments', style: 'currency' },
	long_term_debt: { label: 'Long term debt', style: 'currency' },
	current_debt: { label: 'Current debt', style: 'currency' },
	equity: { label: 'Equity', style: 'currency' },
	operating_cash_flow: { label: 'Operating cash flow', style: 'currency' },
	capital_expenditures: { label: 'Capital expenditures', style: 'currency' },
	free_cash_flow: { label: 'Free cash flow', style: 'currency' },
	number_of_shares: { label: 'Number of shares', style: 'decimal' },
	ppe: { label: 'PPE', style: 'currency' }
};

// Derived financials
export const derivedFinancials: Record<
	string,
	{ label: string; style: Intl.NumberFormatOptions['style'] }
> = {
	eps: { label: 'EPS', style: 'decimal' },
	pe: { label: 'P/E', style: 'decimal' },
	evebit: { label: 'EV/EBIT', style: 'decimal' },
	ps: { label: 'P/S', style: 'decimal' },
	pb: { label: 'P/B', style: 'decimal' }
};
