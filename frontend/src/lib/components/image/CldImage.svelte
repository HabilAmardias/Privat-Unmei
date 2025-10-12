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

<div class="relative overflow-hidden" style={`width: ${width}px; height: ${height}px`}>
	{#if isLoading}
		<div class="absolute inset-0 animate-pulse bg-gray-300"></div>
	{/if}
	<CldImage
		{src}
		{alt}
		{width}
		{height}
		class={`${className} ${isLoading ? 'opacity-0' : 'opacity-100'} transition-opacity duration-500`}
		onload={() => {
			isLoading = false;
		}}
	/>
</div>
