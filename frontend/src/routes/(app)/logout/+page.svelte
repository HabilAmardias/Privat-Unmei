<script lang="ts">
	import { goto } from '$app/navigation';
	import { loadingStore } from '$lib/stores/LoadingStore.svelte';
	import { onMount } from 'svelte';
	onMount(() => {
		if (loadingStore.logOutLoadID) {
			goto('/', { replaceState: true });
		} else {
			goto('/home', { replaceState: true });
		}
	});
</script>
