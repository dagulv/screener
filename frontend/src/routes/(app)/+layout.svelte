<script lang="ts">
	import { renderComponent, renderSnippet } from '$components/shadcn/ui/data-table';
	import Header from '$components/table/header.svelte';

	import Table from '$components/table/table.svelte';
	import type { Screener } from '$lib/api/types.gen.js';
	import type { ColumnDef } from '@tanstack/table-core';
	import { createRawSnippet } from 'svelte';
	import * as Sidebar from '$components/shadcn/ui/sidebar/index.js';
	import AppSidebar from '$components/sidebar.svelte';
	import { flags, formatNumber } from '$lib/constants.js';
	import { derivedFinancials, staticFinancials } from '$lib/screener.js';
	import type { CountryCode, Currency } from '$lib/types.js';

	let { data, children } = $props();

	const sidebar = Sidebar.useSidebar();

	const columns: ColumnDef<Screener>[] = [
		{
			accessorKey: 'countryCode',
			header: () => renderComponent(Header, { class: 'meta' }),
			cell: ({ row }) => {
				const flagSnippet = createRawSnippet<[{ url: string; alt: string }]>((getData) => {
					const data = getData();
					return {
						render: () => `<img class="mt-0.5 h-4 meta" src="${data.url}" alt="${data.alt}">`
					};
				});

				const countryCode = row.getValue('countryCode') as CountryCode;

				return renderSnippet(flagSnippet, { url: flags[countryCode], alt: countryCode });
			},
			enableHiding: false,
			enableSorting: false
		},
		{
			accessorKey: 'name',
			header: ({ column }) => renderComponent(Header, { column, label: 'Name' }),
			cell: ({ row }) => row.getValue('name'),
			enableHiding: false
		},
		{
			accessorKey: 'magicRank',
			header: ({ column }) =>
				renderComponent(Header, { column, label: 'Magic rank', type: 'number' }),
			cell: ({ row }) => {
				const magicRankSnippet = createRawSnippet<[number | null]>((getMagicRank) => {
					let output = '-';
					const raw = getMagicRank();

					if (typeof raw === 'number') {
						output = formatNumber(raw, { style: 'decimal' });
					}

					return {
						render: () => `<div class="text-right font-medium">${output}</div>`
					};
				});

				return renderSnippet(magicRankSnippet, row.getValue('magicRank'));
			}
		},
		{
			accessorKey: 'sector',
			header: ({ column }) => renderComponent(Header, { column, label: 'Sector' }),
			cell: ({ row }) => {
				return renderSnippet(valueSnippet, { type: 'string', value: row.getValue('sector') });
			}
		},
		...Object.entries(derivedFinancials).map(([key, value]) => {
			return {
				accessorKey: key,
				header: ({ column }) =>
					renderComponent(Header, { column, label: value.label, type: 'number' }),
				cell: ({ row }) => {
					return renderSnippet(valueSnippet, {
						type: 'number',
						style: value.style,
						currency: row.original.currency,
						value: row.getValue(key)
					});
				}
			};
		}),
		...Object.entries(staticFinancials).map(([key, value]) => {
			return {
				accessorKey: key,
				header: ({ column }) =>
					renderComponent(Header, { column, label: value.label, type: 'number' }),
				cell: ({ row }) => {
					return renderSnippet(valueSnippet, {
						type: 'number',
						style: value.style,
						currency: row.original.currency,
						value: row.getValue(key)
					});
				}
			};
		})
		// {
		// 	accessorKey: 'amount',
		// 	header: () =>
		// 		renderSnippet(
		// 			createRawSnippet(() => ({
		// 				render: () => `<div class="text-right">Amount</div>`
		// 			}))
		// 		),
		// 	cell: ({ row }) => {
		// 		const amountSnippet = createRawSnippet<[string]>((getAmount) => {
		// 			const amount = Number.parseFloat(getAmount());
		// 			const formatted = new Intl.NumberFormat('en-US', {
		// 				style: 'currency',
		// 				currency: 'USD'
		// 			}).format(amount);
		// 			return {
		// 				render: () => `<div class="text-right font-medium">${formatted}</div>`
		// 			};
		// 		});
		// 		return renderSnippet(amountSnippet, row.getValue('amount'));
		// 	}
		// }
	];

	const outputs = {
		number: (
			value: number | null,
			style?: Intl.NumberFormatOptions['style'],
			currency?: Currency
		) =>
			createRawSnippet(() => {
				return {
					render: () =>
						`<div class="text-right">${typeof value === 'number' ? formatNumber(value, { style, currency }) : '-'}</div>`
				};
			}),
		string: (value: string | null) => value
	} as const;
</script>

{#snippet valueSnippet(params: {
	type: 'number' | 'string';
	style?: Intl.NumberFormatOptions['style'];
	currency?: Currency;
	value: any;
})}
	{@const output = outputs[params.type](params.value, params.style, params.currency)}
	<!-- TODO: Check if output is a snippet and fallback to string instead -->
	{#if typeof output === 'string'}
		<div>{output}</div>
	{:else}
		{@render output?.()}
	{/if}
{/snippet}

{#snippet tableSnippet()}
	<Table
		class="h-[calc(100vh-var(--header-height)-1rem)]! sticky top-[calc(var(--header-height)+0.5rem)]"
		id="screener"
		rowId="companyId"
		{columns}
		data={data.data?.items ?? []}
		count={data.data?.meta.total ?? 0}
		visibleColumns={data.columns ?? new Set()}
	/>
{/snippet}

{#if sidebar.isMobile}
	<AppSidebar children={tableSnippet} />
{/if}

<Sidebar.Inset>
	<main class="flex h-full min-w-0 flex-1 gap-2">
		{#if !sidebar.isMobile}
			{@render tableSnippet()}
		{/if}

		<div class="flex-2">
			{@render children()}
		</div>
	</main>
</Sidebar.Inset>
