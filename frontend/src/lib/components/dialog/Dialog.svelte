<script lang="ts">
	import { Dialog, Separator } from 'bits-ui';
	import { X } from '@lucide/svelte';
	import type { Snippet } from 'svelte';
	type dialogProps = {
		children?: Snippet;
		dialogTitle: Snippet;
		dialogContent: Snippet;
		dialogDescription?: Snippet;
		open?: boolean;
	};
	let {
		children,
		dialogContent,
		dialogDescription,
		dialogTitle,
		open = $bindable()
	}: dialogProps = $props();
</script>

<Dialog.Root bind:open>
	{#if children}
		<Dialog.Trigger
			class="rounded-input text-background shadow-mini hover:bg-dark/95 focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden
	  inline-flex h-12 cursor-pointer items-center justify-center
	  whitespace-nowrap rounded-lg bg-[var(--tertiary-color)] px-[21px] text-[15px] font-semibold text-[var(--secondary-color)] transition-colors hover:text-[var(--primary-color)] focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
		>
			{@render children()}
		</Dialog.Trigger>
	{/if}
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
		/>
		<Dialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] rounded-lg border bg-[var(--primary-color)] p-5 text-[var(--tertiary-color)] sm:max-w-[490px] md:w-full"
		>
			<Dialog.Title
				class="flex w-full items-center justify-center text-lg font-semibold tracking-tight"
			>
				{@render dialogTitle()}
			</Dialog.Title>
			{#if dialogDescription}
				<Separator.Root class="bg-muted -mx-5 mb-6 mt-5 block h-px" />
				<Dialog.Description class="text-foreground-alt text-sm">
					{@render dialogDescription()}
				</Dialog.Description>
			{/if}
			<div class="flex flex-col items-center gap-4 pb-11 pt-7">
				{@render dialogContent()}
			</div>
			<Dialog.Close
				class="focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden absolute right-5 top-5 cursor-pointer rounded-md hover:text-[var(--secondary-color)] focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
			>
				<div>
					<X class="text-foreground size-5" />
					<span class="sr-only">Close</span>
				</div>
			</Dialog.Close>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
