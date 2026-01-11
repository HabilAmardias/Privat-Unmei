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
					class="flex w-full cursor-pointer items-center justify-between py-8 font-bold text-[var(--tertiary-color)] transition-all [&[data-state=open]_svg]:rotate-180"
				>
					{item.header}
					<ChevronDown class="transition-transform duration-200" />
				</Accordion.Trigger>
			</Accordion.Header>
			<Accordion.Content class="accordion-content overflow-hidden">
				<div class="pb-4">
					{@html item.content}
				</div>
			</Accordion.Content>
		</Accordion.Item>
	{/each}
</Accordion.Root>

<style>
	@keyframes accordion-down {
		from {
			height: 0;
		}
		to {
			height: var(--bits-accordion-content-height);
		}
	}

	@keyframes accordion-up {
		from {
			height: var(--bits-accordion-content-height);
		}
		to {
			height: 0;
		}
	}

	:global(.accordion-content[data-state='open']) {
		animation: accordion-down 0.2s ease-out;
	}

	:global(.accordion-content[data-state='closed']) {
		animation: accordion-up 0.2s ease-out;
	}
</style>
