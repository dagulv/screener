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

	const orderby: NonNullable<IterateScreenerData['query']>['orderby'] = q.orderby(
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
				z.literal('ppe'),
				z.literal('eps'),
				z.literal('pe'),
				z.literal('evebit'),
				z.literal('ps'),
				z.literal('pb')
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
			// revenue: q.minmax('revenue')
			revenue: url.searchParams.get('revenue') ?? undefined
		}
	});

	return {
		columns: new Set(['id', 'countryCode', 'name'].concat(visibleColumns)),
		data: data.data
	};
};
