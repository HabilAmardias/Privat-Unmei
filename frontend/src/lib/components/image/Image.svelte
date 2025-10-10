<script lang="ts">
	type imageProps = {
		height?: number;
		width?: number;
		src: string;
		alt?: string;
		className?: string;
	};
	let { height = 32, width = 32, src, alt, className }: imageProps = $props();
	let isLoading = $state<boolean>(true);
</script>

<div class="relative overflow-hidden" style={`width: ${width}px; height: ${height}px`}>
	{#if isLoading}
		<div class="absolute inset-0 animate-pulse bg-gray-300"></div>
	{/if}
	<img
		class={`${className} h-full w-full ${isLoading ? 'opacity-0' : 'opacity-100'} transition-opacity duration-500`}
		{src}
		{alt}
		{height}
		{width}
		onload={() => {
			isLoading = false;
		}}
	/>
</div>
