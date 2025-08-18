<script lang="ts" generics="TData extends { id: string }">
	import Button from '$components/shadcn/ui/button/button.svelte';
	import { ArrowDown, ArrowUp, ArrowUpDownIcon } from '@lucide/svelte';
	import type { Column } from '@tanstack/table-core';

	let {
		column,
		label,
		type = 'string',
		class: className
	}: {
		column?: Column<TData>;
		label?: string;
		type?: 'number' | 'string';
		class?: string;
	} = $props();

	const icon = {
		asc: ArrowUp,
		desc: ArrowDown
	};

	const sort = $derived(column?.getIsSorted());
</script>

{#if column}
	<Button
		variant="ghost"
		size="sm"
		class={[
			'-ml-2 flex  h-6 min-w-max flex-1 items-center justify-start p-0',
			type === 'number' && 'justify-end text-right',
			className
		]}
		onclick={() => column?.toggleSorting(sort === 'asc')}
	>
		<span class="min-w-max">
			{label}
		</span>

		{#if column?.getCanSort()}
			{#if sort}
				{@const Icon = icon[sort]}
				<Icon />
			{:else}
				<ArrowUpDownIcon />
			{/if}
		{/if}
	</Button>
{:else}
	<span
		class={[
			'meta flex h-6 min-w-max flex-1 items-center justify-start p-0',
			type === 'number' && 'text-right',
			className
		]}
	></span>
{/if}
