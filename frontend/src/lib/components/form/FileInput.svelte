<script lang="ts">
	import type { Snippet } from 'svelte';

	type fileInputProps = {
		files?: FileList;
		name: string;
		id: string;
		children: Snippet;
	};
	let { files = $bindable(), name, id, children }: fileInputProps = $props();
</script>

<label for={id} class="cursor-pointer">
	{@render children()}
	<input class="hidden" type="file" bind:files {name} {id} />
</label>
