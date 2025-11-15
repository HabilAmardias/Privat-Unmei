<script lang="ts">
	import type { Snippet } from 'svelte';

	type textAreaProps = {
		children: Snippet;
		value?: string;
		name: string;
		id: string;
		placeholder?: string;
		onBlur?: (e?: FocusEvent & { currentTarget: EventTarget & HTMLTextAreaElement }) => void;
		err?: Error;
	};
	let {
		value = $bindable(),
		name,
		children,
		id,
		placeholder,
		onBlur,
		err
	}: textAreaProps = $props();
</script>

<div>
	<label for={id}>{@render children()}</label>
	{#if err}
		<p class="my-0 text-sm text-[red]">{err.message}</p>
	{/if}
	<textarea
		{id}
		{name}
		onblur={onBlur}
		{placeholder}
		class="w-full rounded-lg bg-[var(--tertiary-color)] text-[var(--secondary-color)]"
		bind:value
	></textarea>
</div>
