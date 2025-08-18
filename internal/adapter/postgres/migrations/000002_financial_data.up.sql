CREATE TABLE financials (
    company_id text NOT NULL REFERENCES companies(id)
        on update cascade
        on delete cascade,
    fiscal_year INT NOT NULL,

    currency text not null references currencies(id),

    -- income statement
    revenue INT,
    cost_of_revenue INT,
    gross_operating_profit INT,
    ebit INT,
    net_income INT,

    -- balance sheet
    total_assets INT,
    total_liabilities INT,
    cash_and_equivalents INT,
    short_term_investments INT,
    long_term_debt INT,
    current_debt INT,
    equity INT,

    -- cash flow
    operating_cash_flow INT,
    capital_expenditures INT,
    free_cash_flow INT,

    -- indexes for fast querying
    primary key (company_id, fiscal_year)
);
