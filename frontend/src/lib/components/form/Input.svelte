<script lang="ts">
	type inputProps = {
		placeholder: string;
		name: string;
		value?: string;
		id: string;
		type: 'text' | 'email';
		onBlur?: (e: FocusEvent & { currentTarget: EventTarget & HTMLInputElement }) => void;
		err?: Error;
	};

	let { placeholder, value = $bindable(), name, onBlur, type, id, err }: inputProps = $props();
</script>

<div class="flex w-full flex-col rounded-md">
	{#if err}
		<p class="my-0 text-sm text-[red]">{err.message}</p>
	{/if}
	<label class="rounded-md bg-[var(--tertiary-color)]" for={id}>
		<input
			bind:value
			onblur={onBlur}
			class="placeholder:text-[var(--secondary-color)]/60 w-full border-none bg-transparent text-[var(--secondary-color)]"
			{id}
			{name}
			{placeholder}
			{type}
		/>
	</label>
</div>
