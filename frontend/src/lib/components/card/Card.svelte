<script lang="ts">
	import type { Snippet } from 'svelte';

	type cardProps = {
		children: Snippet;
	};
	let { children }: cardProps = $props();
</script>

<div class="flex h-fit w-fit flex-col gap-4 rounded-lg bg-none p-3">
	{@render children()}
</div>
