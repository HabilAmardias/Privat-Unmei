<script lang="ts">
	import '../app.css';

	let { children } = $props();
	import { Toaster } from 'svelte-french-toast';
</script>

<div class="app">

	<main class="max-w-[1200px] mx-auto">
		{@render children()}
	</main>
	<Toaster/>
</div>
