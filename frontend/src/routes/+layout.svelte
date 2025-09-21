<script lang="ts">
	import '../app.css';

	let { children } = $props();
	import { Toaster } from 'svelte-french-toast';
</script>

<div class="mx-auto flex h-screen max-h-screen max-w-[1200px] flex-col">
	{@render children()}
	<Toaster />
</div>
