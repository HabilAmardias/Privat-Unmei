<script lang="ts">
	import { Select } from 'bits-ui';
	import { Check, ChevronsUpDown, ChevronsDown, ChevronsUp } from '@lucide/svelte';

	type selectProps = {
		value?: string;
		options: { value: string; label: string }[];
		onChange?: (val: string) => void;
		defaultLable: string;
		name: string;
	};

	let { value = $bindable(), options, onChange, defaultLable, name }: selectProps = $props();

	const selectedLabel = $derived(
		value ? options.find((opt) => opt.value === value)?.label : defaultLable
	);
</script>

<Select.Root
	bind:value
	type="single"
	{name}
	onValueChange={onChange}
	items={options}
	allowDeselect={true}
>
	<Select.Trigger
		class="h-input rounded-9px border-border-input bg-background data-placeholder:text-foreground-alt/50 inline-flex w-full touch-none select-none items-center rounded-lg border border-[var(--tertiary-color)] bg-[var(--tertiary-color)] p-2 px-[11px] text-sm text-[var(--secondary-color)] transition-colors"
	>
		{selectedLabel}
		<ChevronsUpDown
			class="text-muted-foreground ml-auto size-6 cursor-pointer hover:text-[var(--primary-color)]"
		/>
	</Select.Trigger>
	<Select.Portal>
		<Select.Content
			class="focus-override border-muted bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 outline-hidden z-50 h-fit max-h-[var(--bits-select-content-available-height)] w-[var(--bits-select-anchor-width)] min-w-[var(--bits-select-anchor-width)] select-none rounded-xl border bg-[var(--tertiary-color)] px-1 py-3 data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1"
			sideOffset={10}
		>
			<Select.ScrollUpButton class="flex w-full items-center justify-center">
				<ChevronsUp class="size-3 text-[var(--secondary-color)]" />
			</Select.ScrollUpButton>
			<Select.Viewport class="p-1">
				{#each options as opt, i (i + opt.value)}
					<Select.Item
						class="rounded-button data-highlighted:bg-muted outline-hidden data-disabled:opacity-50 data-disabled:text-[var(--secondary-color)] data-disabled:cursor-not-allowed flex h-10 w-full cursor-pointer select-none items-center px-1 py-3 text-sm capitalize text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
						value={opt.value}
						label={opt.label}
					>
						{#snippet children({ selected })}
							{opt.label}
							{#if selected}
								<div class="ml-auto">
									<Check aria-label="check" />
								</div>
							{/if}
						{/snippet}
					</Select.Item>
				{/each}
			</Select.Viewport>
			<Select.ScrollDownButton class="flex w-full items-center justify-center">
				<ChevronsDown class="size-3 text-[var(--secondary-color)]" />
			</Select.ScrollDownButton>
		</Select.Content>
	</Select.Portal>
</Select.Root>
