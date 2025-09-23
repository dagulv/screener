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
		// {
		// 	label: 'Growth',
		// 	value: 'growth'
		// },
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
			label: 'EBIT (MUSD)',
			value: 'ebit',
			categories: ['income-statement'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'Equity (MUSD)',
			value: 'equity',
			categories: ['balance-sheet'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'Gross Operating Profit (MUSD)',
			value: 'gross_operating_profit',
			categories: ['income-statement'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'Net Income (MUSD)',
			value: 'net_income',
			categories: ['income-statement'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'Operating Cash Flow (MUSD)',
			value: 'operating_cash_flow',
			categories: ['cash-flow'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'Revenue (MUSD)',
			value: 'revenue',
			categories: ['income-statement'],
			min: 0,
			max: 1_000_000
		},
		{
			label: 'EPS (USD)',
			value: 'eps',
			categories: ['income-statement'],
			min: 0,
			max: 1_000
		},
		{
			label: 'EV/EBIT',
			value: 'evebit',
			categories: ['valuation'],
			min: 0,
			max: 12
		},
		{
			label: 'P/B',
			value: 'pb',
			categories: ['valuation'],
			min: 0,
			max: 3
		},
		{
			label: 'P/E',
			value: 'pe',
			categories: ['valuation'],
			min: 0,
			max: 25
		},
		{
			label: 'P/S',
			value: 'ps',
			categories: ['valuation'],
			min: 0,
			max: 5
		},
		{
			label: 'Operating Margin',
			value: 'operating_margin',
			categories: ['ratios'],
			min: 0,
			max: 1
		},
		{
			label: 'Net Margin',
			value: 'net_margin',
			categories: ['ratios'],
			min: 0,
			max: 1
		},
		{
			label: 'Return On Equity (MUSD)',
			value: 'roe',
			categories: ['ratios'],
			min: 0,
			max: 1
		},
		{
			label: 'Return On Capital',
			value: 'roc',
			categories: ['ratios'],
			min: 0,
			max: 1
		},
		{
			label: 'Liabilities to Equity',
			value: 'liabilities_to_equity',
			categories: ['ratios'],
			min: 0,
			max: 2
		},
		{
			label: 'Debt to EBIT',
			value: 'debt_to_ebit',
			categories: ['ratios'],
			min: 0,
			max: 3
		},
		{
			label: 'Debt to Assets',
			value: 'debt_to_assets',
			categories: ['ratios'],
			min: 0,
			max: 1
		},
		{
			label: 'Magic Formula',
			value: 'magicFormula',
			categories: ['strategies'],
			min: 1,
			max: 1_000
		}
	];

	let currentCategory = $state('all');
	const currentFilters = $derived(
		filters.filter((f) => currentCategory === 'all' || f.categories.includes(currentCategory))
	);
</script>

<section class="no-scrollbar flex h-full flex-col gap-0 overflow-auto md:flex-row">
	<h2 class="sr-only">Filter</h2>
	<section
		class="h-max w-full border-b p-4 md:sticky md:top-0 md:h-full md:w-auto md:basis-40 md:border-r"
	>
		<h3 class="mb-4 px-1 text-xs font-bold uppercase text-gray-400">Categories</h3>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
		<ul
			class="flex w-full snap-x snap-mandatory items-start gap-1 overflow-x-auto scroll-smooth md:flex-col"
			onclick={(e) => {
				currentCategory =
					(e.target as HTMLElement).closest('button')?.dataset.category ?? categories[0].value;
			}}
		>
			{#each categories as category}
				<li class="contents">
					<button
						class={[
							'w-full flex-1 shrink-0 cursor-pointer snap-start text-nowrap rounded-sm px-2 py-1 text-start text-xs transition hover:bg-gray-100 md:text-wrap',
							currentCategory === category.value && 'bg-gray-100'
						]}
						data-category={category.value}>{category.label}</button
					>
				</li>
			{/each}
		</ul>
	</section>
	<section class="flex-4 flex flex-col p-4">
		<h3 class="mb-4 px-1 text-xs font-bold uppercase text-gray-400">Filter</h3>
		{#if currentFilters.length > 0}
			<ul class=" flex flex-1 flex-col gap-6 px-1">
				{#each currentFilters as filter}
					<li class="contents">
						<Filter {...filter} />
					</li>
				{/each}
			</ul>
		{/if}
	</section>
</section>
