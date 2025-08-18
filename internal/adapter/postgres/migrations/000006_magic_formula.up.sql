CREATE VIEW magic_formula_rankings AS

WITH
	calcs AS (
		SELECT
			f.company_id,
			f.fiscal_year,
			f.ebit::numeric / (f.ppe + f.total_assets - f.total_liabilities) AS roc,
			f.ebit / ((f.number_of_shares * s.average + f.long_term_debt + f.current_debt - f.short_term_investments - f.cash_and_equivalents)::numeric / 1000000) AS yield
		FROM
			financials f
			INNER JOIN shares s ON f.company_id = s.company_id
			AND s.date = TO_DATE((f.fiscal_year + 1) || '-01-02', 'YYYY-MM-DD')
		WHERE
			s.average > 0
	),
	ranks AS (
		SELECT
			c.company_id,
			c.fiscal_year,
			c.roc,
			c.yield,
			rank() OVER (
				PARTITION BY
					c.fiscal_year
				ORDER BY
					c.roc DESC
			) AS roc_rank,
			rank() OVER (
				PARTITION BY
					c.fiscal_year
				ORDER BY
					c.yield DESC
			) AS yield_rank
		FROM
			calcs c
	)
SELECT
	r.company_id,
	r.fiscal_year,
	r.roc,
	r.yield,
	r.roc_rank,
	r.yield_rank,
	rank() OVER (
		PARTITION BY
			r.fiscal_year
		ORDER BY
			r.roc_rank + r.yield_rank ASC
	) AS RANK
FROM
	ranks r