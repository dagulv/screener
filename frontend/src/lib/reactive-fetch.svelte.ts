import { untrack } from 'svelte';

export type Result<T, TError, TInitialValue extends T | undefined> =
	| SuccessResult<T>
	| LoadingResult<ResolveData<T, TInitialValue>>
	| ErrorResult<TError>;

export type SuccessResult<T> = {
	status: 'success';
	data: T;
	error: undefined;
	revalidating: boolean;
};

export type LoadingResult<T> = {
	status: 'loading';
	data: T;
	error: undefined;
	revalidating: boolean;
};

export type ErrorResult<TError> = {
	status: 'error';
	data: undefined;
	error: TError;
	revalidating: boolean;
};

type ResolveData<T, TInitialValue extends T | undefined> = undefined extends TInitialValue
	? T | undefined
	: T;

/**
 * Turns a stateful promise into a SWR-like resource. Revalidation can be done with
 * SvelteKit's `invalidate`, `invalidateAll`, or `goto` functions, depending on the
 * exact use case.
 */
export function fetch<T, TError = Error, TInitialValue extends T | undefined = undefined>(
	promise: () => Promise<T>,
	initialValue: TInitialValue = undefined as TInitialValue
): Result<T, TError, TInitialValue> {
	let data = $state.raw<T | undefined>(initialValue);
	let status = $state<'loading' | 'success' | 'error'>('loading');
	let error = $state.raw<TError | undefined>(undefined);
	let revalidationCount = $state(0);
	const revalidating = $derived(revalidationCount > 0);

	let currentKey: object;

	$effect(() => {
		const key = (currentKey = {});
		untrack(() => revalidationCount++);

		promise()
			.then((resolvedData) => {
				if (key !== currentKey) {
					return;
				}
				data = resolvedData;
				error = undefined;
				status = 'success';
			})
			.catch((err) => {
				if (key !== currentKey) {
					return;
				}
				data = undefined;
				error = err;
				status = 'error';
			})
			.finally(() => {
				revalidationCount--;
			});
	});

	return {
		get status() {
			return status;
		},
		get data() {
			return data;
		},
		get error() {
			return error;
		},
		get revalidating() {
			return revalidating;
		}
	} as Result<T, TError, TInitialValue>;
}
