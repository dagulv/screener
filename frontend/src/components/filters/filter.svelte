<script lang="ts">
	import { page } from '$app/state';
	import { Input } from '$components/shadcn/ui/input/index.js';
	import { Slider } from '$components/shadcn/ui/slider/index.js';
	import { numberFormatter } from '$lib/constants';
	import { debounce, newURLUpdater } from '$lib/utils.svelte';
	import { tick } from 'svelte';

	let {
		label,
		value: key,
		min = 0,
		max = 100
	}: {
		label: string;
		value: string;
		min?: number;
		minLabel?: string;
		max?: number;
		maxLabel?: string;
	} = $props();

	const urlUpdater = newURLUpdater({ page });

	const minimumStepAmount = 10;
	const minimumStep = 0.05;
	const maximumStep = 1000;
	const step = $derived.by(() => {
		const range = max - min;
		let step = Math.max(
			Math.min(Math.round((max - min) / minimumStepAmount), maximumStep),
			minimumStep
		);

		if (step > 1) {
			while (range % step !== 0 && step > minimumStep) {
				step--;
			}
			while (range % step !== 0 && step < maximumStep) {
				step++;
			}
		} else {
			let steps = Math.round(range / step);

			step = range / steps;

			if (step < minimumStep) step = minimumStep;
			if (step > maximumStep) step = maximumStep;
		}
		return step;
	});

	let sliderOverride = true;
	let value = $derived.by(() => {
		const values = page.url.searchParams.get(key)?.split(',') ?? [];
		let minState = undefined;
		let maxState = undefined;

		const minValue = parseFloat(values[0]);
		if (!isNaN(minValue)) minState = minValue;
		const maxValue = parseFloat(values[1]);
		if (!isNaN(maxValue)) maxState = maxValue;

		const state = $state([minState, maxState]);
		return state;
	});
	let sliderValue = $derived.by(() => {
		const state = $state([value[0] ?? min, value[1] ?? max]);
		sliderOverride = true;
		tick().then(() => (sliderOverride = false));
		return state;
	});

	function setFilter(minValue?: number, maxValue?: number) {
		let value = '';

		const validMinValue = typeof minValue === 'number' && minValue >= min;
		const validMaxValue = typeof maxValue === 'number' && maxValue <= max;

		if (validMinValue || validMaxValue) {
			if (validMinValue) {
				value += minValue;
			}
			value += ',';
			if (validMaxValue) {
				value += maxValue;
			}
		}

		urlUpdater.query(key, value);
	}
	const debounceFilter = debounce(setFilter, 300);
</script>

<fieldset>
	<legend class="mb-2 text-xs uppercase text-gray-600">{label}</legend>

	<div class="flex flex-col gap-2">
		<div class="flex w-full items-center justify-between gap-2">
			<span class="text-xs text-gray-600">{numberFormatter.format(min)}</span>
			<span class="text-xs text-gray-600">{numberFormatter.format(max) + '+'}</span>
		</div>
		<Slider
			class="w-full"
			type="multiple"
			value={sliderValue}
			{min}
			{max}
			{step}
			onValueChange={(v) => {
				if (sliderOverride) {
					return;
				}

				const minValue = v[0] > min ? v[0] : undefined;
				const maxValue = v[1] < max ? v[1] : undefined;

				value[0] = minValue;
				value[1] = maxValue;
			}}
			onValueCommit={(v) => {
				const minValue = v[0] !== min ? v[0] : undefined;
				const maxValue = v[1] !== max ? v[1] : undefined;

				value[0] = minValue;
				value[1] = maxValue;

				debounceFilter(value[0], value[1]);
			}}
		/>
		<div class="flex w-full items-stretch justify-between gap-2">
			<Input
				type="number"
				placeholder={numberFormatter.format(min)}
				{min}
				max={value[1] ?? max}
				value={value[0]}
				limit={min}
				class="max-w-xs"
				{step}
				onchange={(v) => {
					if (typeof v === 'number' || typeof v === 'undefined') {
						value[0] = v;
					}

					debounceFilter(value[0], value[1]);
				}}
			/>
			<Input
				type="number"
				placeholder={numberFormatter.format(max) + '+'}
				min={value[0] ?? min}
				{max}
				value={value[1]}
				limit={max}
				class="max-w-xs"
				{step}
				onchange={(v) => {
					if (typeof v === 'number' || typeof v === 'undefined') {
						value[1] = v;
					}

					debounceFilter(value[0], value[1]);
				}}
			/>
		</div>
	</div>
</fieldset>
