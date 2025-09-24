<script lang="ts">
	import { Button } from 'bits-ui';
	import type { Snippet } from 'svelte';
	export type ButtonProps = {
		children: Snippet;
		onClick?: () => void;
		formAction?: string;
		withBg?: boolean;
		full?: boolean;
		type?: 'submit' | 'button' | 'reset' | null;
	};

	let { children, onClick, formAction, type, full = false, withBg = true }: ButtonProps = $props();

	let buttonClass = $state<string>(
		'cursor-pointer rounded-md bg-[var(--tertiary-color)] p-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)]'
	);
	if (!withBg) {
		buttonClass =
			'cursor-pointer rounded-md bg-none p-2 text-[var(--tertiary-color)] hover:text-[var(--primary-color)]';
	}
	if (full) {
		buttonClass += ' w-full';
	}
</script>

<Button.Root {type} formaction={formAction} onclick={onClick} class={buttonClass}>
	{@render children()}
</Button.Root>
