<script lang="ts">
	import { scaleBand } from 'd3-scale';
	import { BarChart, type ChartContextValue } from 'layerchart';
	import * as Chart from '$components/shadcn/ui/chart';
	import { cubicInOut } from 'svelte/easing';
	import type { Currency, ValueLabel } from '$lib/types';
	import { formatNumber } from '$lib/constants';

	let {
		class: className,
		context = $bindable(),
		height,
		data,
		label,
		style = 'decimal',
		currency = 'SEK'
	}: {
		class?: string;
		context?: ChartContextValue;
		height?: number;
		data: ValueLabel<number>[];
		label?: string;
		style?: Intl.NumberFormatOptions['style'];
		currency?: Currency;
	} = $props();

	const chartConfig = {
		value: { label: label, color: 'var(--chart-2)' }
	} satisfies Chart.ChartConfig;
</script>

<Chart.Container {label} config={chartConfig} class={className}>
	<BarChart
		labels={{ offset: 4 }}
		renderContext="svg"
		{data}
		xScale={scaleBand().padding(0.25)}
		x="label"
		y="value"
		series={[{ key: 'value', label: label, color: chartConfig.value.color }]}
		axis="x"
		rule={false}
		props={{
			bars: {
				stroke: 'none',
				radius: 2,
				rounded: 'all',
				// use the height of the chart to animate the bars
				initialY: (context?.height ?? height ?? 0) * 0.8,
				initialHeight: 0,
				motion: {
					y: { type: 'tween', duration: 500, easing: cubicInOut },
					height: { type: 'tween', duration: 500, easing: cubicInOut }
				}
			},
			highlight: { area: { fill: 'none' } },
			labels: {
				format: (d) => {
					return formatNumber(d, { style, currency });
				}
			}
		}}
	>
		{#snippet tooltip()}
			<Chart.Tooltip hideLabel hideIndicator />
		{/snippet}
	</BarChart>
</Chart.Container>
