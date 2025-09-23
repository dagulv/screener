<script lang="ts">
	import SidebarIcon from '@lucide/svelte/icons/sidebar';
	import SearchForm from './search-form.svelte';
	import { Button } from '$components/shadcn/ui/button/index.js';
	import { Separator } from '$components/shadcn/ui/separator/index.js';
	import * as Sidebar from '$components/shadcn/ui/sidebar/index.js';
	import { getCompany } from '$lib/contexts.svelte';

	const sidebar = Sidebar.useSidebar();

	const company = getCompany();
</script>

<header class="bg-background sticky top-0 z-50 -mb-px flex w-full items-center border-b">
	<div class="h-(--header-height) flex w-full items-center gap-2 px-4">
		<div class="contents lg:hidden">
			<Button class="size-8" variant="ghost" size="icon" onclick={sidebar.toggle}>
				<SidebarIcon />
			</Button>
			<Separator orientation="vertical" class="mr-2 h-4" />
		</div>
		<h1 class="hidden lg:block">
			{company.current?.name}
		</h1>
		<!-- <SearchForm class="w-full sm:ml-auto sm:w-auto" /> -->
	</div>
</header>
