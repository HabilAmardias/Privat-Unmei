<script lang="ts">
	import type { Snippet } from 'svelte';

	type cardProps = {
		children: Snippet;
		title: string;
	};
	let { children, title }: cardProps = $props();
</script>

<div class="flex h-fit w-fit flex-col gap-4 rounded-lg bg-[var(--primary-color)] p-3">
	<h2 class="text-center text-2xl font-bold text-[var(--tertiary-color)]">{title}</h2>
	{@render children()}
</div>
