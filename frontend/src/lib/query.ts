import { z } from 'zod';
import { limitKey, orderKey, pageKey, perPage, searchKey, sortKey } from '$lib/constants';
import type { QueryStringParser } from './types';

export const get = (key: string, url: URL) => url.searchParams.get(key) ?? undefined;

export const totalPages = (total: number) => Math.max(1, Math.ceil(total / perPage));

export function queryStringParser(url: URL): QueryStringParser {
	return {
		orderby(union) {
			return union.parse(get(sortKey, url));
		},
		order(preVal) {
			const p = z.union([z.literal('asc'), z.literal('desc')]).default('asc');
			const val = preVal ?? get(orderKey, url);
			return p.parse(val);
		},
		limit() {
			const p = z.coerce.number().min(0).max(100).default(perPage);
			const val = get(limitKey, url);
			return p.parse(val);
		},
		offset() {
			return (this.page() - 1) * this.limit();
		},
		page: () => {
			const val = get(pageKey, url) ?? 1;
			const p = z.coerce.number().min(1);

			return p.parse(val);
		},
		search: () => get(searchKey, url)
	};
}
