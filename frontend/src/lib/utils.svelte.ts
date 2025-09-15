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

class NumberParser {
	private _group: RegExp;
	private _minusSign: RegExp;
	private _decimal: RegExp;
	private _numeral: RegExp;
	private _index: (d: string) => number | undefined;
	constructor(locale: Intl.LocalesArgument) {
		const format = new Intl.NumberFormat(locale);
		const parts = format.formatToParts(-12345.6);
		const numerals = Array.from({ length: 10 }).map((_, i) => format.format(i));
		const index = new Map(numerals.map((d, i) => [d, i]));
		this._minusSign = new RegExp(`[${parts.find((d) => d.type === 'minusSign')?.value}]`);
		this._group = new RegExp(`[${parts.find((d) => d.type === 'group')?.value}]`, 'g');
		this._decimal = new RegExp(`[${parts.find((d) => d.type === 'decimal')?.value}]`);
		this._numeral = new RegExp(`[${numerals.join('')}]`, 'g');
		this._index = (d: string) => index.get(d);
	}
	parse(input: unknown) {
		if (typeof input === 'number' || typeof input === 'bigint') {
			return Number(input);
		} else if (typeof input !== 'string') {
			return 0;
		}
		const DIRECTION_MARK = /\u061c|\u200e/g;
		return +input
			.trim()
			.replace(DIRECTION_MARK, '')
			.replace(this._group, '')
			.replace(this._decimal, '.')
			.replace(this._numeral, this._index)
			.replace(this._minusSign, '-');
	}
}

const numberParser = new NumberParser('en-US');

export function toNumber(input: unknown, isFloat = false): number {
	return numberParser.parse(input);
}

export function debounce<F extends (...args: never[]) => void>(
	func: F,
	delay = 300
): (...args: Parameters<F>) => void {
	let timer: ReturnType<typeof setTimeout>;

	return (...args: Parameters<F>) => {
		clearTimeout(timer);
		timer = setTimeout(() => func(...args), delay);
	};
}

export function watch<T>(
	getter: () => T,
	effectCallback: (value: T | undefined) => (() => void) | undefined
) {
	let previous: T | undefined = undefined;
	$effect(() => {
		const current = getter();
		const cleanup = effectCallback(previous);
		previous = current;

		return cleanup;
	});
}
