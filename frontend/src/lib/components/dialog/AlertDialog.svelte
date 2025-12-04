<script lang="ts">
	import { enhance } from '$app/forms';
	import type { EnhancementArgs, EnhancementReturn } from '$lib/types';
	import { AlertDialog } from 'bits-ui';
	import type { Snippet } from 'svelte';
	type alertDialogProp = {
		open?: boolean;
		action?: string;
		onClick?: (e: MouseEvent & { currentTarget: EventTarget & HTMLButtonElement }) => void;
		enhancement?: (
			args: EnhancementArgs
		) => ({ result, update }: EnhancementReturn) => Promise<void>;
		children: Snippet;
		description: Snippet;
		title: Snippet;
		buttonDisabled?: boolean;
		full?: boolean;
	};
	let {
		open = $bindable(),
		action,
		children,
		description,
		title,
		onClick,
		enhancement,
		buttonDisabled,
		full = false
	}: alertDialogProp = $props();

	let buttonClass = $state<string>(`rounded-input bg-dark text-background shadow-mini
    hover:bg-dark/95 inline-flex h-fit cursor-pointer select-none
    items-center justify-center whitespace-nowrap rounded-lg
	bg-[var(--tertiary-color)] p-2 font-semibold text-[var(--secondary-color)] transition-all hover:text-[var(--primary-color)]`);
	if (full) {
		buttonClass += ' w-full';
	}
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger class={buttonClass} onclick={onClick} disabled={buttonDisabled}>
		{@render children()}
	</AlertDialog.Trigger>
	<AlertDialog.Portal>
		<AlertDialog.Overlay
			class="data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/80"
		/>
		<AlertDialog.Content
			class="rounded-card-lg bg-background shadow-popover data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 outline-hidden fixed left-[50%] top-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border bg-[var(--primary-color)] p-7 text-[var(--tertiary-color)] sm:max-w-lg md:w-full "
		>
			<div class="flex flex-col gap-4 pb-6">
				<AlertDialog.Title class="text-lg font-semibold tracking-tight">
					{@render title()}
				</AlertDialog.Title>
				<AlertDialog.Description class="text-foreground-alt text-sm">
					{@render description()}
				</AlertDialog.Description>
			</div>
			<form
				use:enhance={enhancement}
				method="POST"
				{action}
				class="flex w-full items-center justify-center gap-2"
			>
				<AlertDialog.Cancel
					type="button"
					class="h-input rounded-input bg-muted shadow-mini hover:bg-dark-10 focus-visible:ring-foreground focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full cursor-pointer items-center justify-center py-2 text-[15px] font-medium transition-all hover:text-[var(--secondary-color)] focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
				>
					Cancel
				</AlertDialog.Cancel>
				<AlertDialog.Action
					type="submit"
					formaction={action}
					class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full cursor-pointer items-center justify-center bg-[var(--tertiary-color)] py-2 text-[15px] font-semibold text-[var(--secondary-color)] transition-all hover:text-[var(--primary-color)] focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
				>
					Continue
				</AlertDialog.Action>
			</form>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
