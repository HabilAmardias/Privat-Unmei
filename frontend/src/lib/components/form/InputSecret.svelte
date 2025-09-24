<script lang="ts">
	import { Eye, EyeClosed } from '@lucide/svelte';
	import { Toggle } from 'bits-ui';
	type inputSecretProps = {
		placeholder: string;
		name: string;
		value?: string;
		id: string;
		onChange?: (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => void;
	};

	let { placeholder, value = $bindable(), name, onChange }: inputSecretProps = $props();
	let open = $state<boolean>(false);
</script>

<div class="rounded-md bg-[var(--tertiary-color)]">
	<label class="flex w-fit items-center justify-center gap-1" for="secret">
		<input
			bind:value
			onchange={onChange}
			class="placeholder:text-[var(--secondary-color)]/60 border-none bg-transparent text-[var(--secondary-color)]"
			id="secret"
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
