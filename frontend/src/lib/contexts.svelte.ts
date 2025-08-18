import { getContext, setContext } from 'svelte';
import type { Company, IdAndName } from './api';

// const userKey = Symbol("user");
// /** @param {ReturnType<typeof getUserContext>} user */
// export function setUserContext(user) {
// 	setContext(userKey, user);
// }
// /** @returns {import("$lib/types").User|undefined} */
// export function getUserContext() {
// 	return getContext(userKey);
// }

// const toastKey = Symbol("toast");
// const toast: { current: HTMLElement | undefined } = $state({
// 	current: undefined,
// });
// /** @returns {toast} */
// export function getToast() {
// 	return getContext(toastKey);
// }
// export function setToast() {
// 	return setContext(toastKey, toast);
// }

const sectorKey = Symbol('sector');
export function getCompany(): contextData<Company> {
	return getContext(sectorKey);
}
export function setCompany() {
	setContext(sectorKey, context());
}

type contextData<T> = {
	current: T | undefined;
};

function context<T>(data: T | undefined = undefined): contextData<T> {
	const state = $state({
		current: data
	});
	return state;
}
