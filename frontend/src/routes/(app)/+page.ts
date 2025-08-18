import { getClient } from '$lib/api-client';
import type { PageLoad } from './$types';
import { iterateCompanies } from '$lib/api';
import { error, redirect } from '@sveltejs/kit';

export const load: PageLoad = async ({ fetch }) => {
	const client = getClient(fetch);

	const companies = await iterateCompanies({
		client,
		query: {
			orderBy: 'name',
			limit: 1
		}
	});

	if ((companies.data?.items.length ?? 0) < 1) {
		error(400, 'missing company');
	}

	redirect(307, `/${companies.data?.items[0].id}`);
};
