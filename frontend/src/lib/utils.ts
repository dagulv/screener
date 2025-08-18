import { goto } from '$app/navigation';
import { page } from '$app/state';
import type { Page } from '@sveltejs/kit';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChild<T> = T extends { child?: any } ? Omit<T, 'child'> : T;
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type WithoutChildren<T> = T extends { children?: any } ? Omit<T, 'children'> : T;
export type WithoutChildrenOrChild<T> = WithoutChildren<WithoutChild<T>>;
export type WithElementRef<T, U extends HTMLElement = HTMLElement> = T & { ref?: U | null };
export type WithElementParentRef<T, U extends HTMLElement = HTMLElement> = T & {
	parentRef?: U | null;
};

/**
 * @typedef opt
 * @property {import('@sveltejs/kit').Page} [page]
 * @property {boolean} [replaceState]
 * @property {boolean} [noScroll]
 */

type RequiredOpt = {
	page: Page;
};
type OptionalOpt = {
	replaceState: boolean;
	noScroll: boolean;
	keepParams: boolean;
};

interface URLUpdaterInterface {
	with(opt: RequiredOpt & Partial<OptionalOpt>): URLUpdater;
	url(url: string): void;
	query(key: string | Record<string, unknown>, value?: unknown): void;
}

class URLUpdater implements URLUpdaterInterface {
	#page: Page = page;
	#replaceState: boolean = false;
	#noScroll: boolean = false;
	#keepParams: boolean = false;

	constructor(opt: RequiredOpt & Partial<OptionalOpt>) {
		this.#page = opt.page;
		this.#replaceState = opt.replaceState ?? this.#replaceState;
		this.#noScroll = opt.noScroll ?? this.#noScroll;
		this.#keepParams = opt.keepParams ?? this.#keepParams;
	}

	with(opt: Partial<RequiredOpt & OptionalOpt>): URLUpdater {
		return new URLUpdater({
			page: opt.page ?? this.#page,
			replaceState: opt.replaceState ?? this.#replaceState,
			noScroll: opt.noScroll ?? this.#noScroll,
			keepParams: opt.keepParams ?? this.#keepParams
		});
	}

	url(url: string): void {
		if (this.#keepParams) {
			url += this.#page.url.search;
		}

		goto(url, {
			replaceState: this.#replaceState,
			noScroll: this.#noScroll,
			keepFocus: true
		});
	}

	query(key: string | Record<string, unknown>, val: unknown = null) {
		const searchParams = new URLSearchParams(this.#page.url.search);

		if (key !== 'page' && searchParams.has('page')) {
			searchParams.delete('page');
		}

		if (Array.isArray(key)) {
			if (val === null) {
				for (const k of key) {
					searchParams.delete(k);
				}
			}
		} else if (typeof key === 'object') {
			for (const k in key) {
				const v = Array.isArray(key[k]) ? key[k].join(',') : key[k];

				if ((Array.isArray(v) && v.length) || v) {
					searchParams.set(k, v);
				} else {
					searchParams.delete(k);
				}
			}
		} else {
			if (Array.isArray(val)) {
				val = val.join(',');
			} else if (val instanceof Date) {
				val = val.toISOString();
			} else if (typeof val === 'number') {
				val = `${val}`;
			}

			if (val === null || val === '') {
				searchParams.delete(key);
			} else {
				searchParams.set(key, val);
			}
		}

		return goto('?' + searchParams.toString(), {
			replaceState: this.#replaceState,
			noScroll: this.#noScroll,
			keepFocus: true
		});
	}
}

export function newURLUpdater(opt: RequiredOpt & Partial<OptionalOpt>): URLUpdater {
	return new URLUpdater(opt);
}

export function listErrors(error: unknown) {
	if (error instanceof Error || typeof error === 'string') {
		notice.setNotice(error.message ?? error, { type: 'error' });
		return;
	}

	if (typeof error !== 'object' || !Array.isArray(error.errors)) {
		return;
	}

	// TODO: How do we handle multiple field errors?
	// errors.set(error.errors);
	const msg = error.errors[0].expect ?? error.errors[0].message;
	notice.setNotice(msg, { type: 'error' });
}
