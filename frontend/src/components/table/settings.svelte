<!-- <script>
	import Button from '$components/shadcn/ui/button/button.svelte';
	import * as DropdownMenu from '$components/shadcn/ui/dropdown-menu';
	import { EllipsisIcon } from '@lucide/svelte';
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button variant="ghost" class="size-8 p-0" {...props}>
				<span class="sr-only">Open menu</span>
				<EllipsisIcon />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="end">
		<DropdownMenu.Label>Columns</DropdownMenu.Label>
		<DropdownMenu.Item onclick={() => navigator.clipboard.writeText(payment.id)}>
			Copy payment ID
		</DropdownMenu.Item>
		<DropdownMenu.Separator />
		<DropdownMenu.Item>View customer</DropdownMenu.Item>
		<DropdownMenu.Item>View payment details</DropdownMenu.Item>
	</DropdownMenu.Content>
</DropdownMenu.Root> -->

<script lang="ts" generics="TData extends { id: string }">
	import Settings2Icon from '@lucide/svelte/icons/settings-2';
	import type { Table } from '@tanstack/table-core';
	import { buttonVariants } from '$components/shadcn/ui/button/index.js';
	import * as DropdownMenu from '$components/shadcn/ui/dropdown-menu/index.js';
	import * as Popover from '$components/shadcn/ui/popover/index.js';
	import * as Command from '$components/shadcn/ui/command/index.js';
	import { tick } from 'svelte';
	import { CheckIcon } from '@lucide/svelte';
	import { derivedFinancials, metrics, staticFinancials } from '$lib/screener';

	let { table, onchange }: { table: Table<TData>; onchange?: () => void } = $props();

	let open = $state(false);

	// We want to refocus the trigger button when the user selects
	// an item from the list so users can continue navigating the
	// rest of the form with the keyboard.
	function closeAndFocusTrigger(triggerId: string) {
		open = false;
		tick().then(() => {
			document.getElementById(triggerId)?.focus();
		});
	}
	const triggerId = $props.id();

	const labels = { ...staticFinancials, ...derivedFinancials, ...metrics };
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger
		class={buttonVariants({
			variant: 'outline',
			size: 'sm',
			class: 'ml-auto flex h-6'
		})}
	>
		<Settings2Icon />
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
									closeAndFocusTrigger(triggerId);
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
