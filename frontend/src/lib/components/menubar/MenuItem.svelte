<script lang="ts">
	import { Menubar } from 'bits-ui';
	import type { Snippet } from 'svelte';
	type menuItemProps = {
		children: Snippet;
		separator?: boolean;
	};
	let { children, separator = false }: menuItemProps = $props();
</script>

<Menubar.Menu>
	{@render children()}
	{#if separator}
		<Menubar.Separator />
	{/if}
</Menubar.Menu>
