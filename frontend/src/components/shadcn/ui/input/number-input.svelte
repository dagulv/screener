<script lang="ts">
	import type {
		FormEventHandler,
		HTMLInputAttributes,
		KeyboardEventHandler
	} from 'svelte/elements';
	import { cn, toNumber, type WithElementRef } from '$lib/utils.svelte.js';
	import { numberFormatter } from '$lib/constants';
	import { tick } from 'svelte';

	type Props = WithElementRef<
		Omit<HTMLInputAttributes, 'type' | 'onchange'> & {
			onchange: (value?: number) => void;
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

	const max = $derived(toNumber(restProps.max));
	const min = $derived(toNumber(restProps.min));
	const step = $derived(toNumber(restProps.step));
	let value = $derived.by(() => {
		let state = $state(typeof anyValue === 'undefined' ? limit : toNumber(anyValue));
		return state;
	});
	let valueString = $derived(typeof anyValue === 'undefined' ? '' : numberFormatter.format(value));
	// svelte-ignore state_referenced_locally
	let lastValidValue = $state<number>(getValidValue(value));
	//True when oninput false when onarrows
	let forceBoundary = $state(false);

	function isValid(value: number): boolean {
		return value >= min && value <= max;
	}

	function getValidValue(value: number): number {
		return isValid(value) ? value : (lastValidValue ?? limit);
	}

	function getValueWithinLimit(value: number): number {
		return (limit === min && value < min) || (limit === max && value > max) ? limit : value;
	}

	function onchange(value: number): void {
		if (!isValid(value)) {
			return;
		}
		lastValidValue = value;

		onChange(typeof anyValue === 'undefined' && value === limit ? undefined : value);
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

		switch (e.key) {
			case 'ArrowUp':
				value = getValueWithinLimit(oldValue + step);
				if (value === limit) {
					anyValue = undefined;
				}
				break;
			case 'ArrowDown':
				value = getValueWithinLimit(oldValue - step);
				if (value === limit) {
					anyValue = undefined;
				}
				break;
		}
		console.log(value, oldValue, anyValue);

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

	async function oninput(e: Parameters<FormEventHandler<HTMLInputElement>>[0]) {
		const target = e.target as HTMLInputElement;
		const newValue = parseNumber(target.value);
		const selectionStart = target.selectionStart ?? 0;
		let leftDigits = 0;
		for (let i = 0; i < selectionStart; i++) {
			if (!isNaN(+target.value[i])) leftDigits++;
		}

		try {
			if (newValue !== limit && newValue === value) {
				return;
			}

			value = getValueWithinLimit(newValue);

			onchange(value);
		} finally {
			await tick();

			const newStringValue = numberFormatter.format(value);

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

<div class="flex flex-col justify-start gap-0.5">
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
		<i class="px-2 text-xs text-red-800">{numberFormatter.format(value)} is invalid</i>
	{/if}
</div>
