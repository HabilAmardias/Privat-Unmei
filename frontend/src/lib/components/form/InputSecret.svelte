<script lang="ts">
	import { Eye, EyeClosed } from '@lucide/svelte';
	import { Toggle } from 'bits-ui';
	type inputSecretProps = {
		placeholder: string;
		name: string;
		value?: string;
		id: string;
		onBlur?: (e: FocusEvent & { currentTarget: EventTarget & HTMLInputElement }) => void;
		err?: Error;
	};

	let { placeholder, value = $bindable(), name, onBlur, id, err }: inputSecretProps = $props();
	let open = $state<boolean>(false);
</script>

<div class="flex flex-col rounded-md">
	{#if err}
		<p class="text-sm text-[red]">{err.message}</p>
	{/if}
	<label
		class="flex w-fit items-center justify-center gap-1 rounded-md bg-[var(--tertiary-color)]"
		for="secret"
	>
		<input
			bind:value
			onblur={onBlur}
			class="placeholder:text-[var(--secondary-color)]/60 border-none bg-transparent text-[var(--secondary-color)]"
			{id}
			{name}
			{placeholder}
			type={open ? 'text' : 'password'}
		/>
		<Toggle.Root
			class="pr-2 text-[var(--secondary-color)] hover:text-[var(--primary-color)]"
			bind:pressed={open}
		>
			{#if open}
				<EyeClosed />
			{:else}
				<Eye />
			{/if}
		</Toggle.Root>
	</label>
</div>
