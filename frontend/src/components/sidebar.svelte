<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import * as Sidebar from '$components/shadcn/ui/sidebar/index.js';
	import Button from './shadcn/ui/button/button.svelte';
	import { SidebarIcon } from '@lucide/svelte';
	import Separator from './shadcn/ui/separator/separator.svelte';
	import { getCompany } from '$lib/contexts.svelte';

	let {
		ref = $bindable(null),
		children,
		...restProps
	}: ComponentProps<typeof Sidebar.Root> = $props();

	const sidebar = Sidebar.useSidebar();
	const company = getCompany();
</script>

<Sidebar.Root class="top-(--header-height) h-[calc(100svh-var(--header-height))]!" {...restProps}>
	<Sidebar.Header>
		<Sidebar.MenuButton size="lg" class="max-w-max">
			<Button class="size-8" variant="ghost" size="icon" onclick={sidebar.toggle}>
				<SidebarIcon />
			</Button>
			<Separator orientation="vertical" class="mr-2 h-4" />
		</Sidebar.MenuButton>
		<h1 class="text-xl">
			{company.current?.name}
		</h1>
	</Sidebar.Header>
	<Sidebar.Content class="h-full">
		{@render children?.()}
	</Sidebar.Content>
</Sidebar.Root>
