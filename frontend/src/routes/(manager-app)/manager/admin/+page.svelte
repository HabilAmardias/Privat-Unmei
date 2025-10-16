<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageProps } from './$types';
	import { goto } from '$app/navigation';
	import CldImage from '$lib/components/image/CldImage.svelte';

	let { data }: PageProps = $props();
	onMount(() => {
		if (!data.isVerified) {
			goto('/manager/admin/verify', { replaceState: true });
		}
	});
</script>

<svelte:head>
	<title>My Profile - Privat Unmei</title>
	<meta name="description" content="My Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col gap-4 p-4">
	<div class="flex items-center gap-4">
		<CldImage src={data.profile.profile_image} className="rounded-full" />
	</div>
</div>
