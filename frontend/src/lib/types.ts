/** @typedef {'system_admin' | 'tenant_admin' | 'base' | 'authenticator'} Role */

import type { Icon as IconType } from '@lucide/svelte';
import type {
	HTMLAnchorAttributes,
	HTMLAttributes,
	HTMLButtonAttributes,
	HTMLInputAttributes,
	HTMLLabelAttributes
} from 'svelte/elements';
import type { Snippet } from 'svelte';
import type { Result } from './reactive-fetch.svelte';
import type z from 'zod';

export type Session = {
	id: string;
	loggedIn: boolean;
};

export type User = {
	id: string;
	firstName: string;
	lastName: string;
};

export type ServerError = {
	code: string;
	message: string;
	location: string;
	expect: [string];
};
export type Icon = typeof IconType;
export type ButtonAttributes = HTMLButtonAttributes &
	HTMLAnchorAttributes & {
		variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link';
		size?: 'default' | 'xs' | 'sm' | 'lg' | 'icon';
		icon?: Icon;
		children?: Snippet;
	};
export type TitleAttributes = HTMLAttributes<HTMLHeadingElement> &
	HTMLLabelAttributes & {
		tag?: 'h1' | 'h2' | 'h3' | 'label' | 'span';
		status?: 'normal' | 'error';
	};

export type InputAttributes = HTMLInputAttributes & { error?: ServerError | null };

export type QueryStringParser = {
	orderby: <T extends z.ZodDefault<z.ZodUnion>>(union: T) => z.output<T>;
	order: (preVal?: 'asc' | 'desc') => 'asc' | 'desc';
	limit: () => number;
	offset: () => number;
	page: () => number;
	search: () => string | undefined;
	minmax: (key: string) => string | undefined;
};

export type Response<T> = { rows: Array<T>; meta: { total: number } };
export type TableResult<T> = Result<Response<T>, Error, Response<T>>;
export type ValueLabel<T> = {
	value: T;
	label: string;
};
export type Currency = 'SEK' | 'EUR' | 'DKK' | 'ISK' | 'USD';
export type CountryCode = 'se' | 'dk' | 'fi' | 'is';

export type Filter = {
	label: string;
	value: string;
	categories: string[];
	min: number;
	max: number;
};
