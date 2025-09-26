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
		withPadding?: boolean;
		textColor?: 'light' | 'dark';
		disabled?: boolean;
	};

	let {
		children,
		onClick,
		formAction,
		type,
		full = false,
		withBg = true,
		withPadding = true,
		textColor = 'light',
		disabled
	}: ButtonProps = $props();

	let buttonClass = $state<string>('cursor-pointer rounded-md hover:text-[var(--primary-color)]');
	if (withBg) {
		buttonClass += ' bg-[var(--tertiary-color)]';
	} else {
		buttonClass += ' bg-none';
	}

	if (textColor === 'light') {
		buttonClass += ' text-[var(--secondary-color)]';
	} else {
		buttonClass += ' text-[var(--tertiary-color)]';
	}

	if (full) {
		buttonClass += ' w-full';
	}
	if (withPadding) {
		buttonClass += ` p-2`;
	}
</script>

<Button.Root {disabled} {type} formaction={formAction} onclick={onClick} class={buttonClass}>
	{@render children()}
</Button.Root>
