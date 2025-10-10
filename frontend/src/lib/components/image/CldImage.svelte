<script lang="ts">
	import { CldImage } from 'svelte-cloudinary';
	type cldImageProps = {
		src: string;
		alt?: string;
		width?: number;
		height?: number;
		className?: string;
	};
	let { src, alt, width = 32, height = 32, className }: cldImageProps = $props();

	let isLoading = $state<boolean>(true);
</script>

<CldImage {src} {alt} {width} {height} class={className} onload={() => (isLoading = false)} />
{#if isLoading}
	<div class={`${className} animate-pulse bg-gray-300 w-[${width}px] h-[${height}px]`}></div>
{/if}
