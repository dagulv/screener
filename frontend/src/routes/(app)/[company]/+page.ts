import { getClient } from '$lib/api-client';
import type { PageLoad } from './$types';
import { getCompany, iterateFinancials } from '$lib/api';

export const load: PageLoad = async ({ fetch, params }) => {
	const client = getClient(fetch);

	//TODO: Move await into pages dependent components
	const company = await getCompany({
		client,
		path: { id: params.company }
	});
	const financials = await iterateFinancials({
		client,
		path: { id: params.company }
	});
	return {
		company: company.data,
		financials: financials.data
	};
};
