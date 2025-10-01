<script lang="ts">
	import type { Snippet } from 'svelte';

	type textAreaProps = {
		children: Snippet;
		value?: string;
		name: string;
		id: string;
	};
	let { value = $bindable(), name, children, id }: textAreaProps = $props();
</script>

<div>
	<label for={id}>{@render children()}</label>
	<textarea
		{id}
		{name}
		class="w-full rounded-lg bg-[var(--tertiary-color)] text-[var(--secondary-color)]"
		bind:value
	></textarea>
</div>
