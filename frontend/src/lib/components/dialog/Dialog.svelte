<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Dialog, type WithoutChild } from 'bits-ui';
	import { X } from '@lucide/svelte';

	type Props = Dialog.RootProps & {
		buttonText: string;
		title: Snippet;
		description: Snippet;
		contentProps?: WithoutChild<Dialog.ContentProps>;
		buttonClass?: string;
		// ...other component props if you wish to pass them
	};

	let {
		open = $bindable(false),
		children,
		buttonText,
		contentProps,
		title,
		description,
		...restProps
	}: Props = $props();
</script>

<Dialog.Root bind:open {...restProps}>
	<Dialog.Trigger class={restProps.buttonClass}>
		{buttonText}
	</Dialog.Trigger>
	<Dialog.Portal>
		<Dialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
		/>
		<Dialog.Content
			class="shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 flex w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] flex-col gap-4 rounded-lg border bg-[var(--primary-color)] p-5 sm:max-w-[490px] md:w-full"
			{...contentProps}
		>
			<div class="flex items-center justify-between">
				<Dialog.Title class="text-2xl font-bold text-[var(--tertiary-color)]">
					{@render title()}
				</Dialog.Title>
				<Dialog.Close
					class="cursor-pointer text-[var(--tertiary-color)] hover:text-[var(--secondary-color)]"
					><X /></Dialog.Close
				>
			</div>

			<Dialog.Description>
				{@render description()}
			</Dialog.Description>
			{@render children?.()}
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
