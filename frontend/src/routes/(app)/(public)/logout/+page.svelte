<script lang="ts">
	import { onMount } from 'svelte';
	import { enhance } from '$app/forms';

	onMount(() => {
		document.querySelector('form')?.requestSubmit();
	});
</script>

<form method="POST" use:enhance></form>
