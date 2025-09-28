<script lang="ts">
	import ChartBar from '$components/charts/chart-bar.svelte';
	import type { DerivedFinancialData, FinancialData, Financials } from '$lib/api/types.gen.js';
	import { getCompany } from '$lib/contexts.svelte.js';
	import { derivedFinancials, staticFinancials } from '$lib/screener.js';
	import type { ValueLabel } from '$lib/types.js';

	let { data } = $props();

	const company = getCompany();

	$effect(() => {
		company.current = data.company;
	});

	function formatData(
		data: Financials[],
		key: keyof (FinancialData & DerivedFinancialData)
	): ValueLabel<number>[] {
		const formattedData = data
			.map((f) => {
				return {
					label: String(f.fiscalYear),
					value:
						f.staticData?.[key as keyof FinancialData] ??
						f.derivedData?.[key as keyof DerivedFinancialData] ??
						0
				};
			})
			.sort((a, b) => (a.label > b.label ? 1 : -1));

		return formattedData;
	}

	const properties = [
		...Object.entries(staticFinancials),
		...Object.entries(derivedFinancials)
	].map(([key, meta]) => {
		return {
			value: key,
			label: meta.label
		};
	});
</script>

<h1 class="mb-8 mt-4 block text-center text-xl lg:hidden">
	{company.current?.name}
</h1>

<div class="grid grid-cols-[repeat(auto-fill,minmax(14rem,1fr))] gap-2">
	{#each properties as property}
		<ChartBar
			label={`${property.label} (${data.company?.currency?.name})`}
			class="h-40 w-full rounded-sm border border-gray-200 bg-gray-50 p-2"
			data={formatData(data.financials?.items ?? [], property.value as keyof FinancialData)}
			style="decimal"
		/>
	{/each}
</div>
