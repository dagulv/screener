import { iterateScreener, type IterateScreenerData } from '$lib/api';
import { getClient } from '$lib/api-client';
import { screenerKey, sortKey } from '$lib/constants';
import z from 'zod';
import type { LayoutLoad } from '../$types';
import { queryStringParser } from '$lib/query';
import { newURLUpdater } from '$lib/utils.svelte';
import { page } from '$app/state';

const defaultVisibleColumns = ['magicRank'];

export const load: LayoutLoad = async ({ fetch, url, depends }) => {
	depends('app:screener');

	const client = getClient(fetch);

	const q = queryStringParser(url);
	let visibleColumns = localStorage.getItem(`table-columns/${screenerKey}`)?.split(',');

	const orderby = q.orderby(
		z
			.union([
				z.literal('name'),
				z.literal('magicRank'),
				z.literal('sector'),
				z.literal('revenue'),
				z.literal('cost_of_revenue'),
				z.literal('gross_operating_profit'),
				z.literal('ebit'),
				z.literal('net_income'),
				z.literal('total_assets'),
				z.literal('total_liabilities'),
				z.literal('cash_and_equivalents'),
				z.literal('short_term_investments'),
				z.literal('long_term_debt'),
				z.literal('current_debt'),
				z.literal('equity'),
				z.literal('operating_cash_flow'),
				z.literal('capital_expenditures'),
				z.literal('free_cash_flow'),
				z.literal('number_of_shares'),
				z.literal('eps'),
				z.literal('evebit'),
				z.literal('pb'),
				z.literal('pe'),
				z.literal('ps'),
				z.literal('operating_margin'),
				z.literal('net_margin'),
				z.literal('roe'),
				z.literal('roc'),
				z.literal('liabilities_to_equity'),
				z.literal('debt_to_ebit'),
				z.literal('debt_to_assets'),
				z.literal('cash_conversion')
			])
			.default('name')
	);

	if (visibleColumns) {
		if (!visibleColumns.includes(orderby)) {
			return newURLUpdater({ page }).query(sortKey, 'name');
		}
	} else {
		visibleColumns = defaultVisibleColumns;
	}

	//TODO: Move await into table component
	const data = await iterateScreener({
		client,
		query: {
			orderby: orderby,
			order: q.order(),
			offset: q.offset(),
			limit: q.limit(),
			search: q.search(),
			include: url.searchParams.get('include')?.split(',') ?? undefined,
			columns: visibleColumns,
			capital_expenditures: q.minmax('capital_expenditures'),
			ebit: q.minmax('ebit'),
			equity: q.minmax('equity'),
			gross_operating_profit: q.minmax('gross_operating_profit'),
			net_income: q.minmax('net_income'),
			operating_cash_flow: q.minmax('operating_cash_flow'),
			revenue: q.minmax('revenue'),
			eps: q.minmax('eps'),
			evebit: q.minmax('evebit'),
			pb: q.minmax('pb'),
			pe: q.minmax('pe'),
			ps: q.minmax('ps'),
			operating_margin: q.minmax('operating_margin'),
			net_margin: q.minmax('net_margin'),
			roe: q.minmax('roe'),
			roc: q.minmax('roc'),
			liabilities_to_equity: q.minmax('liabilities_to_equity'),
			debt_to_ebit: q.minmax('debt_to_ebit'),
			debt_to_assets: q.minmax('debt_to_assets'),
			cash_conversion: q.minmax('cash_conversion'),
			magicRank: q.minmax('magicRank')
		}
	});

	return {
		columns: new Set(['id', 'countryCode', 'name'].concat(visibleColumns)),
		data: data.data
	};
};
