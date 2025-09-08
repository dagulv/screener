import type { CountryCode, Currency } from './types';
import seFlag from 'circle-flags/flags/se.svg';
import dkFlag from 'circle-flags/flags/dk.svg';
import fiFlag from 'circle-flags/flags/fi.svg';
import isFlag from 'circle-flags/flags/is.svg';

export const perPage = 50;
export const screenerKey = 'screener';
export const orderKey = 'order';
export const limitKey = 'limit';
export const pageKey = 'page';
export const searchKey = 's';
export const sortKey = 'sort';
const numberFormatLocale: Intl.LocalesArgument = 'se-SV';
const numberFormatOptions: Intl.NumberFormatOptions = {
	maximumFractionDigits: 2
};
const numberPrefixes = { B: 1_000_000_000, M: 1_000_000, K: 1_000 };
const currencies: Currency[] = ['SEK', 'EUR', 'DKK', 'ISK'];
export const numberFormatter = new Intl.NumberFormat(numberFormatLocale);
export const flags: Record<CountryCode, string> = {
	se: seFlag,
	dk: dkFlag,
	fi: fiFlag,
	is: isFlag
};
export const formatters = {
	currency: currencies.reduce(
		(acc, curr) => {
			acc[curr] = new Intl.NumberFormat(numberFormatLocale, {
				currency: curr,
				...numberFormatOptions
			});
			return acc;
		},
		{} as Record<Currency, Intl.NumberFormat>
	),
	percent: new Intl.NumberFormat(numberFormatLocale, { style: 'percent', ...numberFormatOptions }),
	decimal: new Intl.NumberFormat(numberFormatLocale, { style: 'decimal', ...numberFormatOptions })
};

const defaultOpts: { style: Intl.NumberFormatOptions['style']; currency?: Currency } = {
	style: 'currency',
	currency: 'SEK'
};
export function formatNumber(
	number: number,
	newOpts?: { style: Intl.NumberFormatOptions['style']; currency?: Currency }
): string {
	const opts = { ...defaultOpts, ...newOpts };

	let prefix = '';

	for (const [key, value] of Object.entries(numberPrefixes)) {
		if (Math.abs(number) >= value) {
			number /= value;
			prefix = key;
			break;
		}
	}

	if (opts.style === 'currency') {
		if (!opts.currency) {
			throw Error('currency is required on currency style');
		}

		return formatters.currency[opts.currency].format(number) + ' ' + prefix + opts.currency;
	}

	return formatters[opts.style].format(number) + ' ' + prefix;
}
