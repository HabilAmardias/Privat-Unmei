<script lang="ts">
	import { Menubar } from 'bits-ui';
	import type { Snippet } from 'svelte';
	type menubarProps = {
		children: Snippet;
	};
	let { children }: menubarProps = $props();
</script>

<Menubar.Root
	class="fixed bottom-0 left-0 right-0 mx-auto flex w-full justify-between bg-[var(--tertiary-color)] px-3 py-2 md:bottom-auto md:top-0 md:my-1 md:w-[60%] md:rounded-lg"
>
	{@render children()}
</Menubar.Root>
