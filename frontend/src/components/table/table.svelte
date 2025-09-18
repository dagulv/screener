<script lang="ts" generics="TData extends { [rowId]: string }">
	import { createSvelteTable } from '$components/shadcn/ui/data-table';
	import {
		getCoreRowModel,
		type ColumnDef,
		type ColumnFiltersState,
		type PaginationState,
		type Row,
		type RowSelectionState,
		type SortingState,
		type VisibilityState
	} from '@tanstack/table-core';
	import * as Table from '$components/shadcn/ui/table';
	import Button from '$components/shadcn/ui/button/button.svelte';
	import FlexRender from '$components/shadcn/ui/data-table/flex-render.svelte';
	import { ChevronLeft, ChevronRight, Settings2 } from '@lucide/svelte';
	import { orderKey, pageKey, perPage, sortKey } from '$lib/constants';
	import Settings from './settings.svelte';
	import { newURLUpdater } from '$lib/utils.svelte';
	import { page } from '$app/state';
	import { invalidate } from '$app/navigation';
	import { tick } from 'svelte';
	import * as Sheet from '$components/shadcn/ui/sheet/index.js';
	import Filters from '$components/filters/filters.svelte';

	let {
		id,
		rowId = 'id',
		columns: rawColumns,
		data,
		count,
		selectedRows = $bindable([]),
		visibleColumns,
		order = { orderby: 'name', order: 'asc' },
		class: className
	}: {
		id: string;
		rowId: string;
		columns: ColumnDef<TData>[];
		data: TData[];
		count: number;
		selectedRows?: string[];
		visibleColumns: Set<string>;
		order?: { orderby: string; order: 'asc' | 'desc' };
		class: string;
	} = $props();

	if (page.url.searchParams.has(sortKey) || page.url.searchParams.has(orderKey)) {
		const orderValue = page.url.searchParams.get(orderKey);
		order = {
			orderby: page.url.searchParams.get(sortKey) ?? order.orderby,
			order: orderValue === 'asc' || orderValue === 'desc' ? orderValue : order.order
		};
	}

	const urlUpdater = newURLUpdater({ page });

	const columns: ColumnDef<TData>[] = [
		// TODO: Add checkbox when its needed.
		// {
		// 	id: 'select',
		// 	header: ({ table }) =>
		// 		renderComponent(Checkbox, {
		// 			checked: table.getIsAllPageRowsSelected(),
		// 			indeterminate: table.getIsSomeRowsSelected() && !table.getIsAllRowsSelected(),
		// 			onCheckedChange: (v) => table.toggleAllPageRowsSelected(!!v),
		// 			'aria-label': 'Select all',
		// 			class: 'mt-0'
		// 		}),
		// 	cell: ({ row }) =>
		// 		renderComponent(Checkbox, {
		// 			checked: row.getIsSelected(),
		// 			onCheckedChange: (v) => row.toggleSelected(!!v),
		// 			class: 'mt-0.5'
		// 		}),
		// 	enableSorting: false,
		// 	enableHiding: false
		// },
		...rawColumns
	];

	async function onchangeColumnToggle() {
		if (!tableRef) {
			return;
		}

		await tick();

		tableRef.scrollLeft = tableRef.scrollWidth;
	}

	const initialColumnVisibility = {
		...rawColumns.reduce((acc: VisibilityState, curr) => {
			if (curr.accessorKey) {
				acc[curr.accessorKey] = false;
			}
			return acc;
		}, {}),
		...Array.from(visibleColumns).reduce((acc: VisibilityState, curr) => {
			if (curr) {
				acc[curr] = true;
			}
			return acc;
		}, {})
	};

	let tableRef = $state<HTMLElement | null>(null);

	let pagination = $state<PaginationState>({
		pageIndex: +(page.url.searchParams.get(pageKey) ?? 1) - 1,
		pageSize: perPage
	});
	let sorting = $state<SortingState>([{ id: order.orderby, desc: order.order === 'desc' }]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let rowSelection = $state<RowSelectionState>({});
	let columnVisibility = $state<VisibilityState>(initialColumnVisibility);

	const table = createSvelteTable({
		get data() {
			return data;
		},
		manualFiltering: true,
		manualSorting: true,
		manualPagination: true,
		get rowCount() {
			return count;
		},
		columns,
		state: {
			get pagination() {
				return pagination;
			},
			get sorting() {
				return sorting;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			},
			get columnFilters() {
				return columnFilters;
			}
		},
		getCoreRowModel: getCoreRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}

			urlUpdater.query(pageKey, pagination.pageIndex + 1);
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}

			table.firstPage();
			urlUpdater.query({
				[sortKey]: sorting[0].id,
				[orderKey]: sorting[0].desc ? 'desc' : 'asc'
			});
		},
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}

			//Filter rows here
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}

			let cols: string[] = [];

			for (const [key, value] of Object.entries(columnVisibility)) {
				if (value) {
					cols.push(key);
				}
			}
			localStorage.setItem(`table-columns/${id}`, cols.join(','));
			invalidate('app:screener');
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}

			selectedRows = table.getSelectedRowModel().rows.map((r) => r.original[rowId]);
		}
	});

	const tableState = $derived(table.getState());

	function onclickRow(row: Row<TData>) {
		urlUpdater.with({ keepParams: true }).url(`/${row.original.companyId}`);
	}
</script>

<div
	class="sticky top-[calc(var(--header-height)+0.5rem)] flex h-full min-w-0 max-w-full flex-1 flex-col [--table-filter-header-height:calc(var(--spacing)*8)]"
>
	<div
		class="ml-auto flex h-[var(--table-filter-header-height)] min-h-[var(--table-filter-header-height)] items-center gap-1 py-1"
	>
		<Sheet.Root>
			<Sheet.Trigger>
				{#snippet child({ props })}
					<Button variant="outline" size="icon" class="h-full" {...props}>
						<Settings2 />
					</Button>
				{/snippet}
			</Sheet.Trigger>
			<Sheet.Content class="sm:max-w-lg">
				<Filters />
			</Sheet.Content>
		</Sheet.Root>

		<Settings class="h-full" {table} onchange={onchangeColumnToggle} />
	</div>
	<section
		class={[
			'rounded-md border [--footer-height:calc(var(--spacing)*8)] [--table-header-height:calc(var(--spacing)*8)]',
			className
		]}
		style={`--column-count: ${visibleColumns.size - 1};`}
	>
		<Table.Root bind:parentRef={tableRef}>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row>
						{#each headerGroup.headers as header (header.id)}
							<Table.Head data-name={header.id}>
								{#if !header.isPlaceholder}
									<FlexRender
										context={header.getContext()}
										content={header.column.columnDef.header}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#if table.getRowModel().rows?.length}
					{#each table.getRowModel().rows as row (row.id)}
						<Table.Row
							data-state={row.getIsSelected() && 'selected'}
							onclick={() => onclickRow(row)}
						>
							{#each row.getVisibleCells() as cell (cell.id)}
								<Table.Cell data-name={cell.column.id}>
									<FlexRender context={cell.getContext()} content={cell.column.columnDef.cell} />
								</Table.Cell>
							{/each}
						</Table.Row>
					{/each}
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
					</Table.Row>
				{/if}
			</Table.Body>
		</Table.Root>
		<div class="flex h-[var(--footer-height)] items-center justify-end gap-2 border-t px-3 py-1">
			<div class="text-muted-foreground h-full flex-1 content-center text-xs">
				{table.getRowCount()} total row(s)
			</div>
			<div class="flex h-full items-center gap-2">
				<div class="text-muted-foreground flex-1 content-center text-xs">
					{tableState.pagination.pageIndex * tableState.pagination.pageSize + 1}
					-
					{(tableState.pagination.pageIndex + 1) * tableState.pagination.pageSize}
				</div>
				<Button
					variant="outline"
					size="sm"
					class="h-full"
					onclick={() => table.previousPage()}
					disabled={!table.getCanPreviousPage()}
				>
					<ChevronLeft />
				</Button>
				<Button
					variant="outline"
					size="sm"
					class="h-full"
					onclick={() => table.nextPage()}
					disabled={!table.getCanNextPage()}
				>
					<ChevronRight />
				</Button>
			</div>
		</div>
	</section>
</div>

<style>
	@reference "tailwindcss";

	section {
		:global {
			table {
				@apply grid h-full w-full select-none select-auto auto-rows-auto grid-cols-[repeat(var(--column-count),1fr)];

				thead {
					@apply contents;

					tr th {
						@apply border-b border-gray-100;
					}
				}

				th {
					@apply py-1;
				}

				tbody {
					vertical-align: baseline;
					@apply relative z-[1] contents overflow-y-scroll;

					tr:not(:last-of-type) td {
						@apply border-b border-gray-100;
					}
				}

				tr {
					@apply contents;

					td,
					th {
						@apply m-0 flex-[1_1_0] align-top;

						p {
							@apply leading-4;
						}

						.meta {
							@apply w-4;
						}
					}

					th {
						@apply flex items-center justify-start;
					}

					td {
						* {
							@apply min-w-0 overflow-hidden overflow-ellipsis;
						}

						&.misc {
							@apply p-0;
						}
					}
				}
			}
		}
	}
</style>
