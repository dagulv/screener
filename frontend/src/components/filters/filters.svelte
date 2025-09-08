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
			label: 'Revenue',
			value: 'revenue',
			categories: ['income-statement'],
			min: 0,
			minLabel: '0 MUSD',
			max: 1_000_000,
			maxLabel: '1 000 000 MUSD'
		}
	];

	let currentCategory = $state('all');
	const currentFilters = $derived(
		filters.filter((f) => currentCategory === 'all' || f.categories.includes(currentCategory))
	);
</script>

<section>
	<h2 class="sr-only">Filter</h2>
	<section>
		<h3>Kategorier</h3>
		<ul
			onchange={(e) =>
				(currentCategory =
					(e.target as HTMLElement).closest('button')?.dataset.category ?? categories[0].value)}
		>
			{#each categories as category}
				<li>
					<button data-category={category.value}>{category.label}</button>
				</li>
			{/each}
		</ul>
	</section>
	<section>
		<h3>Filter</h3>
		{#if currentFilters.length > 0}
			<ul>
				{#each currentFilters as filter}
					<Filter {...filter} />
				{/each}
			</ul>
		{/if}
	</section>
</section>
