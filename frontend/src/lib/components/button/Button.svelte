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

	let buttonClass = $state<string>(
		'flex gap-2 items-center justify-center cursor-pointer rounded-md hover:text-[var(--primary-color)] disabled:cursor-not-allowed disabled:opacity-50'
	);
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
	} else {
		buttonClass += ' w-fit';
	}
	if (withPadding) {
		buttonClass += ` p-2`;
	}
</script>

<Button.Root {disabled} {type} formaction={formAction} onclick={onClick} class={buttonClass}>
	{@render children()}
</Button.Root>
