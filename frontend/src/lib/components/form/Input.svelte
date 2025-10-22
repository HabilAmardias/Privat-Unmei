<script lang="ts">
	type inputProps = {
		placeholder?: string;
		name: string;
		value?: string;
		id: string;
		type: 'text' | 'email' | 'number' | 'time';
		onBlur?: (e: FocusEvent & { currentTarget: EventTarget & HTMLInputElement }) => void;
		err?: Error;
		width?: 'full' | number;
		min?: number;
		max?: number;
		step?: number;
	};

	let {
		placeholder,
		value = $bindable(),
		name,
		onBlur,
		type,
		id,
		err,
		width = 'full',
		min,
		max,
		step
	}: inputProps = $props();
	let containerClass = $state<string>('flex flex-col rounded-md');
	if (width === 'full') {
		containerClass += ` w-full`;
	} else {
		containerClass += ` w-[${width}px]`;
	}
</script>

<div class={containerClass}>
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
			{min}
			{max}
			{step}
		/>
	</label>
</div>
