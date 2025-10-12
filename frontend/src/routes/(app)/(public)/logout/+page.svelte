<script lang="ts">
	import { goto } from '$app/navigation';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { onMount } from 'svelte';
	onMount(() => {
		if (loadingStore.logOutLoadID) {
			goto('/login', { replaceState: true });
		} else {
			goto('/', { replaceState: true });
		}
	});
</script>
