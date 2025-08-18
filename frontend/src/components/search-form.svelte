<script lang="ts">
	import type { HTMLFormAttributes } from 'svelte/elements';
	import { newURLUpdater, type WithElementRef } from '$lib/utils.js';
	import { iterateCompanies, type Company, type IdAndName } from '$lib/api';
	import Autocomplete from './autocomplete.svelte';
	import { page } from '$app/state';
	let { ref = $bindable(null), ...restProps }: WithElementRef<HTMLFormAttributes> = $props();
	const urlUpdater = newURLUpdater({ page, keepParams: true });

	async function getCompanies(search: string): Promise<Company & { id: string; name: string }[]> {
		const data = await iterateCompanies({ query: { search } });

		return (data.data?.items ?? []) as Company & { id: string; name: string }[];
	}

	function onchange(item: IdAndName) {
		urlUpdater.url(`/${item.id}`);
	}
</script>

<!-- <Command.Item
							value={subitem.name}
							onSelect={() => {
								if (values.has(/** @type {string} */ (subitem.id))) {
									values.delete(/** @type {string} */ (subitem.id));
								} else {
									values.add(/** @type {string} */ (subitem.id));
								}
								updateQueryParam({ page }, /** @type {string} */ (item.value), Array.from(values));
							}}
						>
							{subitem.name}
							<CheckIcon
								class={cn(!values.has(/** @type {string} */ (subitem.id)) && 'text-transparent')}
							/>
						</Command.Item> -->

<!-- {#snippet itemSnippet(item: Company)}
	test
{/snippet} -->

<form {...restProps} bind:this={ref}>
	<Autocomplete
		items={getCompanies}
		emptyMessage="No results."
		placeholder="Search companies"
		onValueChange={onchange}
	/>
</form>
