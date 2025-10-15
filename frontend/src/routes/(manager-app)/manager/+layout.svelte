<script lang="ts">
	import { onMount } from 'svelte';
	import type { LayoutProps } from './$types';
	import { adminRole, mentorRole } from '$lib/utils/constants';
	import { goto } from '$app/navigation';

	const { children, data }: LayoutProps = $props();

	onMount(() => {
		if (data.role !== mentorRole && data.role !== adminRole) {
			goto('/courses', { replaceState: true });
		}
	});
</script>

{@render children()}
