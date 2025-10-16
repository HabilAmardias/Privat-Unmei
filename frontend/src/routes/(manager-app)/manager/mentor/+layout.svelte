<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { adminRole, mentorRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';

	const { children, data }: LayoutProps = $props();

	onMount(() => {
		if (data.role !== mentorRole) {
			goto('/courses', { replaceState: true });
		}
	});
</script>

<main class="h-screen">
	{@render children()}
</main>
