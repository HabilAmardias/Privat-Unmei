<script lang="ts">
	import { page } from '$app/state';
	import CldImage from '$lib/components/image/CldImage.svelte';
	import { ErrorLogo } from '$lib/utils/constants';
</script>

<div class="flex h-dvh w-full items-center justify-center">
	<div class="flex flex-col items-center">
		<CldImage src={ErrorLogo} />
		<h1 class="font-bold text-[var(--tertiary-color)]">{page.status}</h1>
		<p class="text-[var(--tertiary-color)]">{page.error?.message}</p>
	</div>
</div>
