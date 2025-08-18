<script lang="ts">
	import '../app.css';
	import * as Sidebar from '$components/shadcn/ui/sidebar/index.js';
	import SiteHeader from '$components/site-header.svelte';
	import { setCompany } from '$lib/contexts.svelte';

	let { children } = $props();

	setCompany();
</script>

<div class="[--header-height:calc(--spacing(14))]">
	<Sidebar.Provider class="flex flex-col">
		<SiteHeader />
		<div>
			{@render children()}
		</div>
	</Sidebar.Provider>
</div>
