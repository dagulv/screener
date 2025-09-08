<script lang="ts">
	import type {
		HTMLInputAttributes,
		HTMLInputTypeAttribute,
		KeyboardEventHandler
	} from 'svelte/elements';
	import { cn, toNumber, type WithElementRef } from '$lib/utils.js';
	import { numberFormatter } from '$lib/constants';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'onchange'> & {
			onchange?: (value: number) => void;
			limit?: number;
		}
	>;

	let {
		ref = $bindable(null),
		value: anyValue = $bindable(),
		files = $bindable(),
		class: className,
		limit,
		onchange,
		...restProps
	}: Props = $props();

	const max = toNumber(restProps.max);
	const min = toNumber(restProps.min);
	const step = toNumber(restProps.step);
	let value = $derived.by(() => {
		const state = $state(toNumber(anyValue));
		return state;
	});

	function onkeydown(e: Parameters<KeyboardEventHandler<HTMLInputElement>>[0]) {
		const inputValue = (e.target as HTMLInputElement).value;

		let newValue = 0;
		// console.log(inputValue || 0);

		switch (e.key) {
			case 'ArrowUp':
				newValue = Number(inputValue || 0) + step;
				if (value <= max) {
					value = newValue;
				}
			case 'ArrowDown':
				newValue = Number(inputValue || 0) - step;
				if (value >= min) {
					value = newValue;
				}
		}

		// onchange?.(value);
	}

	function parseNumber(v: any): void {
		value = toNumber(v);
		console.log(v, value);
	}

	function formatNumber(): string {
		if (value !== limit) {
			return numberFormatter.format(value);
		}

		return '';
	}
</script>

<input
	bind:this={ref}
	data-slot="input"
	class={cn(
		'border-input bg-background selection:bg-primary dark:bg-input/30 selection:text-primary-foreground ring-offset-background placeholder:text-muted-foreground shadow-xs flex h-9 w-full min-w-0 rounded-md border px-3 py-1 text-base outline-none transition-[color,box-shadow] disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
		'focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px]',
		'aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive',
		className
	)}
	type="text"
	pattern="(?:0|[1-9]\d*)"
	inputMode="decimal"
	autoComplete="off"
	{onkeydown}
	bind:value={formatNumber, parseNumber}
	{...restProps}
/>
