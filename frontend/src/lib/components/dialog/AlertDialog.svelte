<script lang="ts">
	import { AlertDialog } from 'bits-ui';
	import type { Snippet } from 'svelte';
	type alertDialogProp = {
		open: boolean;
		action?: string;
        children: Snippet
        description: Snippet
        title: Snippet
		onSubmit?: (e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }) => void;
	};
	let { open = $bindable(), action, onSubmit, children, description, title }: alertDialogProp = $props();
</script>

<AlertDialog.Root bind:open>
	<AlertDialog.Trigger
		class="rounded-input bg-dark text-background shadow-mini
    hover:bg-dark/95 inline-flex h-12 cursor-pointer select-none items-center
    justify-center whitespace-nowrap rounded-lg bg-[var(--tertiary-color)]
	p-2 px-[21px] text-[15px] font-semibold text-[var(--secondary-color)] transition-all hover:text-[var(--primary-color)] active:scale-[0.98]"
	>
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
				onsubmit={onSubmit}
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
					type={onSubmit ? 'submit' : 'button'}
					class="h-input rounded-input bg-dark text-background shadow-mini hover:bg-dark/95 focus-visible:ring-dark focus-visible:ring-offset-background focus-visible:outline-hidden inline-flex w-full cursor-pointer items-center justify-center bg-[var(--tertiary-color)] py-2 text-[15px] font-semibold text-[var(--secondary-color)] transition-all hover:text-[var(--primary-color)] focus-visible:ring-2 focus-visible:ring-offset-2 active:scale-[0.98]"
				>
					Continue
				</AlertDialog.Action>
			</form>
		</AlertDialog.Content>
	</AlertDialog.Portal>
</AlertDialog.Root>
