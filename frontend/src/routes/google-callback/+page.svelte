<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import Loading from '$lib/components/loader/Loading.svelte';
	let { data }: PageProps = $props();
	onMount(() => {
		if (data.success) {
			goto('/courses', { replaceState: true });
		} else {
			goto('/login', { replaceState: true });
		}
	});
</script>

<div class="h-dvh">
	<Loading />
</div>
