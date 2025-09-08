<script lang="ts">
	import { page } from '$app/state';
	import { Input } from '$components/shadcn/ui/input/index.js';
	import { Slider } from '$components/shadcn/ui/slider/index.js';
	import { formatNumber, numberFormatter } from '$lib/constants';

	let {
		label,
		value: key,
		min = 0,
		minLabel = '',
		max = 100,
		maxLabel = ''
	}: {
		label: string;
		value: string;
		min?: number;
		minLabel?: string;
		max?: number;
		maxLabel?: string;
	} = $props();

	const step = Math.max(Math.round((max - min) / 1000), 1);

	let { min: minValue, max: maxValue } = $derived.by(() => {
		const values = page.url.searchParams.get(key)?.split(',') ?? [];
		let minState = $state(min);
		let maxState = $state(max);

		const minValue = parseInt(values[0]);
		if (!isNaN(minValue)) minState = minValue;
		const maxValue = parseInt(values[1]);
		if (!isNaN(maxValue)) maxState = maxValue;

		return {
			get min() {
				return minState;
			},
			get max() {
				console.log('maxState', maxState);

				return maxState;
			}
		};
	});
</script>

<fieldset>
	<legend>{label}</legend>

	<div class="flex flex-col gap-1">
		<div class="flex w-full items-center justify-between gap-2">
			<span class="text-xs text-gray-600">{minLabel}</span>
			<span class="text-xs text-gray-600">{maxLabel}</span>
		</div>
		<Slider
			class="w-full"
			type="multiple"
			bind:value={
				() => [minValue, maxValue],
				(value) => {
					minValue = value[0];
					maxValue = value[1];
				}
			}
			{min}
			{max}
			{step}
			onValueCommit={(value) => console.log(value)}
		/>
		<div class="flex w-full items-center justify-between gap-2">
			<Input
				type="number"
				placeholder={String(min)}
				{min}
				max={maxValue}
				bind:value={minValue}
				class="max-w-xs"
				{step}
				onchange={() => console.log(minValue)}
			/>
			<Input
				type="number"
				placeholder={String(max)}
				min={minValue}
				{max}
				limit={max}
				bind:value={
					() => numberFormatter.format(maxValue),
					(value) => {
						maxValue = +value;
						console.log(+value, maxValue);
					}
				}
				class="max-w-xs"
				{step}
				onchange={() => console.log(maxValue)}
			/>
		</div>
	</div>
</fieldset>
