<script lang="ts">
	import type { Snippet } from 'svelte';

	type fileInputProps = {
		files?: FileList;
		name: string;
		id: string;
		children: Snippet;
		accept?: string;
	};
	let { files = $bindable(), name, id, children, accept }: fileInputProps = $props();
</script>

<label for={id} class="cursor-pointer">
	{@render children()}
	<input {accept} class="hidden" type="file" bind:files {name} {id} />
</label>
