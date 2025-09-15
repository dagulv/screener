<script lang="ts">
	import type {
		FormEventHandler,
		HTMLInputAttributes,
		KeyboardEventHandler
	} from 'svelte/elements';
	import { cn, toNumber, type WithElementRef } from '$lib/utils.svelte.js';
	import { numberFormatter } from '$lib/constants';
	import { tick } from 'svelte';
	import { de } from 'zod/locales';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'onchange'> & {
			onchange: (value: number) => void;
			limit?: number;
		}
	>;

	let {
		ref = $bindable(null),
		value: anyValue,
		files = $bindable(),
		class: className,
		limit = 0,
		onchange: onChange,
		...restProps
	}: Props = $props();

	const max = toNumber(restProps.max);
	const min = toNumber(restProps.min);
	const step = toNumber(restProps.step);
	let value = $derived.by(() => {
		let state = $state(toNumber(anyValue));
		return state;
	});
	let valueString = $derived(formatNumber());
	// svelte-ignore state_referenced_locally
	let lastValidValue = $state<number>(getValidValue(value));

	function isValid(value: number): boolean {
		return value >= min && value <= max;
	}

	function getValidValue(value: number): number {
		return isValid(value) ? value : (lastValidValue ?? limit);
	}

	function onchange(value: number): void {
		if (!isValid(value)) {
			return;
		}
		lastValidValue = value;

		onChange(value);
	}

	function onkeydown(e: Parameters<KeyboardEventHandler<HTMLInputElement>>[0]) {
		switch (e.key) {
			case 'ArrowUp':
			case 'ArrowDown':
				handleArrows(e);
				break;
			case 'ArrowLeft':
			case 'ArrowRight':
				break;
		}
	}

	async function handleArrows(e: Parameters<KeyboardEventHandler<HTMLInputElement>>[0]) {
		e.preventDefault();
		let oldValue = parseNumber((e.target as HTMLInputElement).value);
		let newValue = 0;

		switch (e.key) {
			case 'ArrowUp':
				newValue = oldValue + step;
				break;
			case 'ArrowDown':
				newValue = oldValue - step;
				break;
		}

		if (newValue >= min && newValue <= max) {
			value = newValue;
		} else {
			value = limit;
		}

		onchange(value);

		await tick();
		const end = (ref as HTMLInputElement)?.value.length;
		(ref as HTMLInputElement)?.setSelectionRange(end, end);
		ref?.focus();
	}

	function parseNumber(v: string): number {
		if (v.trim() === '') {
			return limit;
		}

		const current = toNumber(v);

		if (isNaN(current)) {
			return value;
		}

		return current;
	}

	function formatNumber(): string {
		if (value !== limit) {
			return numberFormatter.format(value);
		}

		return '';
	}

	async function oninput(e: Parameters<FormEventHandler<HTMLInputElement>>[0]) {
		const target = e.target as HTMLInputElement;
		const newValue = parseNumber(target.value);
		const selectionStart = target.selectionStart ?? 0;
		let leftDigits = 0;
		for (let i = 0; i < selectionStart; i++) {
			if (!isNaN(+target.value[i])) leftDigits++;
		}

		try {
			if (newValue === value) {
				return;
			}

			value = newValue;

			onchange(value);
		} finally {
			await tick();

			const newStringValue = formatNumber();

			(e.target as HTMLInputElement).value = newStringValue;

			let digitsSeen = 0;
			let newPos = 0;
			while (digitsSeen < leftDigits && newPos < newStringValue.length) {
				if (!isNaN(+newStringValue[newPos])) digitsSeen++;
				newPos++;
			}

			(e.target as HTMLInputElement).setSelectionRange(newPos, newPos);
		}
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
	pattern="\d*"
	inputMode="decimal"
	autoComplete="off"
	{oninput}
	{onkeydown}
	value={valueString}
	{...{ ...restProps, ...{ min: undefined, max: undefined, step: undefined } }}
/>

{#if !isValid(value)}
	{value} is invalid
{/if}
