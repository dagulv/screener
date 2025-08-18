<script lang="ts" module>
	export type Getter<TData extends Data> = (search: string) => Promise<TData[]>;
	type Data =
		| { label: string | number; id: string | number }
		| { label: string | number; value: string | number }
		| { name: string | number; id: string | number }
		| { name: string | number; value: string | number };

	function getLabel<TData extends Data>(item: TData) {
		if ('label' in item) return String(item.label);
		if ('name' in item) return String(item.name);
		throw new Error('label not found');
	}

	function getValue<TData extends Data>(item: TData) {
		if ('id' in item) return String(item.id);
		if ('value' in item) return String(item.value);
		throw new Error('value not found');
	}

	function arrayGetter<TData extends Data>(items: TData[]): Getter<TData> {
		return (search: string) => {
			search = search.toLowerCase();
			return new Promise((resolve) => {
				resolve(items.filter((item) => getLabel(item).toLowerCase().indexOf(search) !== -1));
			});
		};
	}
</script>

<script lang="ts" generics="TData extends Data">
	import * as Command from '$components/shadcn/ui/command/index.js';
	import { cn } from '$lib/utils';
	import { Check } from '@lucide/svelte';
	import type { FormEventHandler } from 'svelte/elements';
	import Skeleton from './shadcn/ui/skeleton/skeleton.svelte';

	type AutoCompleteProps = {
		items: TData[] | Getter<TData>;
		emptyMessage: string;
		value?: TData;
		onValueChange?: (value: TData) => void;
		disabled?: boolean;
		placeholder?: string;
	};

	let {
		items: rawItems,
		placeholder,
		emptyMessage,
		value,
		onValueChange,
		disabled
	}: AutoCompleteProps = $props();

	let inputRef = $state<HTMLInputElement>(null!);

	let isOpen = $state(false);
	let selected = $state<TData | undefined>(value);
	let inputValue = $state<string>(value ? getLabel(value) : '');

	const getter = $derived(Array.isArray(rawItems) ? arrayGetter(rawItems) : rawItems);
	let data = $state<Promise<TData[]> | undefined>(undefined);
	$effect(() => {
		data = getter(inputValue);
	});
	const oninput: FormEventHandler<HTMLInputElement> = (e) => {
		if (!isOpen) {
			isOpen = true;
		}

		switch (e.key) {
			case 'Enter':
				if (inputRef.value !== '') {
					const items = getter(inputRef.value);
					items.then((items) => {
						// const optionToSelect = items.find((item) => getLabel(item) === inputRef.value);
						const optionToSelect = items[0];
						if (optionToSelect) {
							selected = optionToSelect;
							onValueChange?.(optionToSelect);
						}
					});
				}
			case 'Escape':
				inputRef.blur();
		}
	};

	function handleBlur() {
		isOpen = false;
	}

	function handleSelectOption(selectedOption: TData) {
		inputValue = '';

		selected = selectedOption;

		onValueChange?.(selectedOption);

		// This is a hack to prevent the input from being focused after the user selects an option
		// We can call this hack: "The next tick"
		setTimeout(() => {
			inputRef?.blur();
		}, 0);
	}
</script>

<Command.Root class="overflow-visible" shouldFilter={false}>
	<Command.Input
		bind:ref={inputRef}
		bind:value={inputValue}
		onblur={handleBlur}
		onfocus={() => (isOpen = true)}
		{placeholder}
		{disabled}
		class="text-base outline-none"
		{oninput}
		wrapperClass="ring ring-gray-300"
	/>
	{#if isOpen}
		<div class="relative mt-1">
			<div
				class={cn(
					'animate-in fade-in-0 zoom-in-95 absolute top-0 z-10 w-full rounded-xl bg-white outline-none'
				)}
			>
				<Command.List class="max-h-[300px] overflow-y-auto rounded-lg ring-1 ring-slate-200">
					{#await data}
						<Command.Loading>
							<div class="p-1">
								<Skeleton class="h-8 w-full" />
							</div>
						</Command.Loading>
					{:then items}
						{#if items && items.length > 0}
							<Command.Group>
								{#each items as item}
									{@const isSelected = selected && getValue(selected) === getValue(item)}
									<!-- key={item.value} -->
									<Command.Item
										value={getLabel(item)}
										onmousedown={(event) => {
											event.preventDefault();
											event.stopPropagation();
										}}
										onSelect={() => handleSelectOption(item)}
										class={cn('flex w-full items-center gap-2', !isSelected ? 'pl-8' : null)}
									>
										{#if isSelected}
											<Check class="w-4" />
										{/if}
										{getLabel(item)}
									</Command.Item>
								{/each}
							</Command.Group>
						{:else}
							<Command.Empty class="select-none rounded-sm px-2 py-3 text-center text-sm">
								{emptyMessage}
							</Command.Empty>
						{/if}
					{/await}
				</Command.List>
			</div>
		</div>
	{/if}
</Command.Root>
