<script lang="ts">
	import { Toaster } from 'svelte-french-toast';
	import '../app.css';
	import { navigating } from '$app/state';

	let { children } = $props();
</script>

<div class="mx-auto flex max-h-dvh max-w-[1200px] flex-col">
	{#if navigating.to}
		<div class="fixed top-0 left-0 h-2 w-full bg-[var(--tertiary-color)] animate-pulse z-50"></div>
	{/if}
	{@render children()}
	<Toaster />
</div>
