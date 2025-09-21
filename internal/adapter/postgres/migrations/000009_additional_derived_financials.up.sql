drop view derived_financials;
create view derived_financials as

with calcs as (
	select
		f.company_id,
		f.fiscal_year,
		f.net_income::float * 10000 / nullif(f.number_of_shares, 0) as eps,
		f.number_of_shares * s.average::float / 1000000 / nullif(f.net_income, 0) as pe,
		(f.number_of_shares * s.average::float / 1000000 + (f.long_term_debt + f.current_debt - f.short_term_investments - f.cash_and_equivalents)::float) / nullif(f.ebit, 0) AS evebit,
		f.number_of_shares * s.average::float / 1000000 / nullif(f.revenue, 0) as ps,
		f.number_of_shares * s.average::float / 1000000 / nullif((f.total_assets-f.total_liabilities), 0) as pb,

        f.ebit::float / nullif(f.revenue, 0)::float as operating_margin,
        f.ebit::float / nullif(f.net_income, 0)::float as net_margin,
        f.net_income::float / nullif(f.equity, 0)::float as roe,
        f.ebit::float / nullif(f.ppe + f.total_assets - f.total_liabilities, 0)::float as roc,
        f.total_liabilities::float / nullif(f.equity, 0)::float as liabilities_to_equity,
        (f.long_term_debt + f.current_debt - f.cash_and_equivalents)::float / nullif(f.ebit, 0)::float as debt_to_ebit,
        (f.long_term_debt + f.current_debt)::float / nullif(f.total_assets, 0)::float as debt_to_assets,
        f.operating_cash_flow::float / nullif(f.net_income, 0)::float as cash_conversion
	from financials f
	
	inner join shares s on s.company_id = f.company_id and s.date = TO_DATE((f.fiscal_year + 1) || '-01-02', 'YYYY-MM-DD')
)
select 
    c.company_id, 
    c.fiscal_year, 
    c.eps, 
    c.pe, 
    c.evebit, 
    c.ps, 
    c.pb,
    c.operating_margin,
    c.net_margin,
    c.roe,
    c.roc,
    c.liabilities_to_equity,
    c.debt_to_ebit,
    c.debt_to_assets,
    c.cash_conversion
from calcs c