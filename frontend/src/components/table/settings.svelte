<script lang="ts" generics="TData extends { id: string }">
	import SettingsIcon from '@lucide/svelte/icons/settings';
	import type { Table } from '@tanstack/table-core';
	import { buttonVariants } from '$components/shadcn/ui/button/index.js';
	import * as DropdownMenu from '$components/shadcn/ui/dropdown-menu/index.js';
	import * as Command from '$components/shadcn/ui/command/index.js';
	import { CheckIcon } from '@lucide/svelte';
	import { derivedFinancials, metrics, staticFinancials } from '$lib/screener';

	let {
		table,
		onchange,
		class: className
	}: { table: Table<TData>; onchange?: () => void; class?: string } = $props();

	const triggerId = $props.id();

	const labels = { ...staticFinancials, ...derivedFinancials, ...metrics };
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={buttonVariants({
			variant: 'outline',
			size: 'sm',
			class: ['ml-auto flex', className]
		})}
	>
		<SettingsIcon />
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<Command.Root>
			<Command.Input id={triggerId} placeholder="Toggle columns..." />
			<DropdownMenu.Separator />
			<DropdownMenu.Group>
				<Command.List>
					<Command.Empty>No results found.</Command.Empty>
					<Command.Group>
						{#each table
							.getAllColumns()
							.filter((col) => typeof col.accessorFn !== 'undefined' && col.getCanHide()) as column (column)}
							<Command.Item
								value={labels[column.id]?.label ?? column.id}
								onSelect={() => {
									column.toggleVisibility();
									onchange?.();
								}}
							>
								{labels[column.id]?.label ?? column.id}

								<CheckIcon class={['ml-auto', !column.getIsVisible() && 'text-transparent']} />
							</Command.Item>
						{/each}
					</Command.Group>
				</Command.List>
			</DropdownMenu.Group>
		</Command.Root>
	</DropdownMenu.Content>
</DropdownMenu.Root>
