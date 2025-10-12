<script lang="ts">
	import type { Snippet } from 'svelte';

	type cardProps = {
		children: Snippet;
		withBg?: boolean;
		theme?: 'light' | 'dark';
		center?: boolean;
	};
	let { children, withBg = false, theme = 'dark', center = false }: cardProps = $props();

	let containerClass = $state<string>('flex h-fit w-fit flex-col gap-4 rounded-lg p-3');
	if (withBg) {
		if (theme === 'light') {
			containerClass += ' bg-[var(--primary-color)] text-[var(--tertiary-color)]';
		} else {
			containerClass += ' bg-[var(--tertiary-color)]';
		}
	}
	if (center) {
		containerClass += ' justify-center items-center';
	}
</script>

<div class={containerClass}>
	{@render children()}
</div>
