drop view derived_financials;
create view derived_financials as

with calcs as (
	select
		f.company_id,
		f.fiscal_year,
		(f.net_income::float * 10000 / nullif(f.number_of_shares, 0)) as eps,
		(f.number_of_shares * s.average::float / 1000000 / nullif(f.net_income, 0)) as pe,
		((f.number_of_shares * s.average::float / 1000000 + (f.long_term_debt + f.current_debt - f.short_term_investments - f.cash_and_equivalents)::float) / nullif(f.ebit, 0)) AS evebit,
		(f.number_of_shares * s.average::float / 1000000 / nullif(f.revenue, 0)) as ps,
		(f.number_of_shares * s.average::float / 1000000 / nullif((f.total_assets-f.total_liabilities), 0)) as pb
	from financials f
	
	inner join shares s on s.company_id = f.company_id and s.date = TO_DATE((f.fiscal_year + 1) || '-01-02', 'YYYY-MM-DD')
)
select c.company_id, c.fiscal_year, c.eps, c.pe, c.evebit, c.ps, c.pb from calcs c