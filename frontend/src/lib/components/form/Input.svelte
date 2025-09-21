<script lang="ts">
	type inputProps = {
		placeholder: string;
		name: string;
		value?: string;
		id: string;
		type: 'text' | 'email';
		onChange?: (e: Event & { currentTarget: EventTarget & HTMLInputElement }) => void;
	};

	let { placeholder, value = $bindable(), name, onChange, type, id }: inputProps = $props();
</script>

<div class="w-full rounded-md bg-[var(--tertiary-color)]">
	<label class="rounded-md" for={id}>
		<input
			bind:value
			onchange={onChange}
			class="placeholder:text-[var(--secondary-color)]/60 w-full border-none bg-transparent text-[var(--secondary-color)]"
			{id}
			{name}
			{placeholder}
			{type}
		/>
	</label>
</div>
