<script lang="ts">
	import { ChevronDown } from '@lucide/svelte';
	import { Accordion } from 'bits-ui';

	type accordionItems = {
		header: string;
		content: string;
	};

	type accordionProps = {
		type: 'single' | 'multiple';
		items: accordionItems[];
	};

	let { type, items }: accordionProps = $props();
</script>

<Accordion.Root {type}>
	{#each items as item}
		<Accordion.Item class="group border-b border-[var(--tertiary-color)] px-2">
			<Accordion.Header>
				<Accordion.Trigger
					class="flex w-full items-center justify-between py-8 font-bold text-[var(--tertiary-color)]"
				>
					{item.header}
					<ChevronDown />
				</Accordion.Trigger>
			</Accordion.Header>
			<Accordion.Content>
				<div>
					{item.content}
				</div>
			</Accordion.Content>
		</Accordion.Item>
	{/each}
</Accordion.Root>

<style>
</style>
