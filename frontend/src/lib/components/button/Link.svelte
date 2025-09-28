<script lang="ts">
	import { Button } from 'bits-ui';
	import type { Snippet } from 'svelte';
	type anchorProps = {
		href: string;
		children: Snippet;
		theme?: 'light' | 'dark';
	};
	let { href, children, theme = 'light' }: anchorProps = $props();

	let className = $state<string>('text-[var(--secondary-color)] hover:text-[var(--primary-color)]');
	if (theme === 'dark') {
		className = 'text-[var(--tertiary-color)] hover:text-[var(--primary-color)]';
	}
</script>

<Button.Root {href} class={className}>
	{@render children()}
</Button.Root>
