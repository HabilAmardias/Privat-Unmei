<script lang="ts">
	import type { PageProps } from './$types';
	import CldImage from '$lib/components/image/CldImage.svelte';

	let { data }: PageProps = $props();
</script>

<svelte:head>
	<title>Profile - Privat Unmei</title>
	<meta name="description" content="Profile - Privat Unmei" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
</svelte:head>

<div class="flex h-full flex-col gap-4 p-4">
	<div class="flex items-center gap-4">
		<CldImage src={data.profile.profile_image} className="rounded-full" />
		<div class="flex flex-col gap-1">
			<div class="flex gap-1">
				<b class="text-xl text-[var(--tertiary-color)]">{data.profile.name}</b>
			</div>
			<p class="text-md">{data.profile.email}</p>
		</div>
	</div>
	<div class="flex flex-col gap-2">
		<b class="text-xl text-[var(--tertiary-color)]">Bio:</b>
		<p class="text-justify">{data.profile.bio}</p>
	</div>
</div>
