<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	let { data }: PageProps = $props();
	onMount(async () => {
		if (data.success) {
			await goto('/courses', { replaceState: true });
		} else {
			await goto('/login', { replaceState: true });
		}
	});
</script>
