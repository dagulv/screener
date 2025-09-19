<script lang="ts">
	import Filter from './filter.svelte';

	const categories = [
		{
			label: 'All',
			value: 'all'
		},
		{
			label: 'Ratios',
			value: 'ratios'
		},
		{
			label: 'Growth',
			value: 'growth'
		},
		{
			label: 'Valuation',
			value: 'valuation'
		},
		{
			label: 'Income statement',
			value: 'income-statement'
		},
		{
			label: 'Balance sheet',
			value: 'balance-sheet'
		},
		{
			label: 'Cash flow',
			value: 'cash-flow'
		},
		{
			label: 'Strategies',
			value: 'strategies'
		}
	];

	const filters = [
		{
			label: 'Revenue (MUSD)',
			value: 'revenue',
			categories: ['income-statement'],
			min: 0,
			max: 1_000_000
		}
	];

	let currentCategory = $state('all');
	const currentFilters = $derived(
		filters.filter((f) => currentCategory === 'all' || f.categories.includes(currentCategory))
	);
</script>

<section class="flex h-full gap-4">
	<h2 class="sr-only">Filter</h2>
	<section class="h-full border-r p-4">
		<h3 class="mb-4 text-xs font-bold uppercase text-gray-400">Categories</h3>
		<ul
			class="flex flex-col items-start"
			onchange={(e) =>
				(currentCategory =
					(e.target as HTMLElement).closest('button')?.dataset.category ?? categories[0].value)}
		>
			{#each categories as category}
				<li class="contents">
					<button
						class="w-full flex-1 cursor-pointer rounded-sm px-2 py-1 text-start text-xs transition hover:bg-gray-100"
						data-category={category.value}>{category.label}</button
					>
				</li>
			{/each}
		</ul>
	</section>
	<section class="h-full p-4">
		<h3 class="mb-4 text-xs font-bold uppercase text-gray-400">Filter</h3>
		{#if currentFilters.length > 0}
			<ul>
				{#each currentFilters as filter}
					<Filter {...filter} />
				{/each}
			</ul>
		{/if}
	</section>
</section>
