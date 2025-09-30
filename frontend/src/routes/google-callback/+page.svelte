<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	let { data }: PageProps = $props();
	onMount(async () => {
		if (data.success) {
			localStorage.setItem('status', data.status);
			await goto('/home', { replaceState: true });
		}
		await goto('/', { replaceState: true });
	});
</script>
