<script lang="ts">
	import type { HTMLTableAttributes } from 'svelte/elements';
	import { cn, type WithElementParentRef, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		parentRef = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementParentRef<WithElementRef<HTMLTableAttributes>> = $props();
</script>

<div
	bind:this={parentRef}
	data-slot="table-container"
	class="relative h-[calc(100%-var(--table-header-height))] w-full overflow-x-auto"
>
	<table
		bind:this={ref}
		data-slot="table"
		class={cn('w-max caption-bottom text-sm', className)}
		{...restProps}
	>
		{@render children?.()}
	</table>
</div>
