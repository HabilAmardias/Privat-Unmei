<script lang="ts">
	import { Combobox } from 'bits-ui';
	import { ChevronsUp, ChevronsDown } from '@lucide/svelte';
	type searchProps = {
		keyword?: string;
		onKeywordChange?: (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => void;
		onValueChange?: (val: string) => void;
		label: string;
		items?: { value: string; label: string }[];
		value?: string;
	};
	let {
		keyword = '',
		onKeywordChange,
		items,
		onValueChange,
		label,
		value = $bindable()
	}: searchProps = $props();
</script>

<Combobox.Root
	type="single"
	onOpenChangeComplete={(o) => {
		if (!o) keyword = '';
	}}
	bind:value
	{onValueChange}
>
	<div class="relative w-fit rounded-lg bg-[var(--tertiary-color)]">
		<Combobox.Input
			oninput={onKeywordChange}
			class="h-input border-border-input bg-background focus:ring-foreground focus:ring-offset-background focus:outline-hidden inline-flex w-fit touch-none truncate rounded-lg border bg-[var(--tertiary-color)] text-base text-[var(--secondary-color)] transition-colors placeholder:text-[var(--secondary-color)] focus:ring-2 focus:ring-offset-2 sm:text-sm md:w-[296px]"
			placeholder={label}
			aria-label={label}
		/>
	</div>
	{#if items}
		<Combobox.Portal>
			<Combobox.Content
				class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-fit max-h-[var(--bits-combobox-content-available-height)] w-[var(--bits-combobox-anchor-width)] min-w-[var(--bits-combobox-anchor-width)] select-none rounded-xl border bg-[var(--tertiary-color)] px-1 py-3 text-[var(--secondary-color)] data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
				sideOffset={10}
			>
				<Combobox.ScrollUpButton class="flex w-full items-center justify-center py-1">
					<ChevronsUp class="size-3" />
				</Combobox.ScrollUpButton>
				<Combobox.Viewport class="p-1">
					{#each items as item, i (i + item.value)}
						<Combobox.Item
							class="rounded-button data-highlighted:bg-muted outline-hidden flex h-10 w-full select-none items-center py-3 pl-5 pr-1.5 text-sm capitalize hover:cursor-pointer hover:text-[var(--primary-color)]"
							value={item.value}
							label={item.label}
						>
							{#snippet children()}
								{item.label}
							{/snippet}
						</Combobox.Item>
					{:else}
						<span class="block px-5 py-2 text-sm text-[var(--secondary-color)]">
							No results found, try again.
						</span>
					{/each}
				</Combobox.Viewport>
				<Combobox.ScrollDownButton class="flex w-full items-center justify-center py-1">
					<ChevronsDown class="size-3" />
				</Combobox.ScrollDownButton>
			</Combobox.Content>
		</Combobox.Portal>
	{/if}
</Combobox.Root>
