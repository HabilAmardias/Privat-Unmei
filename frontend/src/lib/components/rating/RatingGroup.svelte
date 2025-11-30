<script lang="ts">
	import { unstable_RatingGroup } from 'bits-ui';
	import { Star } from '@lucide/svelte';
	type ratingGroupProps = {
		value: number;
		name: string;
	};
	let { value = $bindable(), name }: ratingGroupProps = $props();
</script>

<unstable_RatingGroup.Root {name} bind:value min={1} max={5} class="flex gap-1">
	{#snippet children({ items })}
		{#each items as item (item.index)}
			<unstable_RatingGroup.Item
				index={item.index}
				class="group size-10 cursor-pointer text-[var(--primary-color)] transition-colors hover:text-[var(--primary-color)] data-[state=inactive]:text-[var(--tertiary-color)] md:size-8"
			>
				<Star class="size-full fill-current" />
			</unstable_RatingGroup.Item>
		{/each}
	{/snippet}
</unstable_RatingGroup.Root>
